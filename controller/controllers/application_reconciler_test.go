package controllers

import (
	"context"
	"github.com/kapp-staging/kapp/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func generateEmptyApplication() *v1alpha1.Application {
	name := randomName()[:12]

	application := &v1alpha1.Application{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Application",
			APIVersion: "core.kapp.dev/v1alpha1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      name,
			Namespace: "test",
		},
		Spec: v1alpha1.ApplicationSpec{
			Components: []v1alpha1.ComponentSpec{},
			SharedEnv:  []v1alpha1.EnvVar{},
		},
	}

	return application
}

const timeout = time.Second * 20
const interval = time.Millisecond * 500

var _ = Describe("Application basic CRUD", func() {
	defer GinkgoRecover()

	It("Should handle application correctly", func() {
		By("Create")
		application := generateEmptyApplication()
		Expect(k8sClient.Create(context.Background(), application)).Should(Succeed())

		By("Ready")
		fetched := &v1alpha1.Application{}
		Eventually(func() error {
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, fetched)
		}, timeout, interval).Should(Succeed())

		By("Update")
		fetched.Spec.SharedEnv = append(fetched.Spec.SharedEnv, v1alpha1.EnvVar{
			Name:  "name",
			Value: "value",
			Type:  v1alpha1.EnvVarTypeStatic,
		})
		Expect(k8sClient.Update(context.Background(), fetched)).Should(Succeed())
		fetchedUpdated := &v1alpha1.Application{}
		Eventually(func() bool {
			_ = k8sClient.Get(context.Background(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, fetchedUpdated)
			return len(fetchedUpdated.Spec.SharedEnv) == 1 && fetchedUpdated.Spec.SharedEnv[0].Value == "value"
		}, timeout, interval).Should(Equal(true))

		By("Delete")
		Eventually(func() error {
			f := &v1alpha1.Application{}
			_ = k8sClient.Get(context.Background(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, f)
			return k8sClient.Delete(context.Background(), f)
		}, timeout, interval).Should(Succeed())

		By("Read after delete")
		Eventually(func() error {
			f := &v1alpha1.Application{}
			return k8sClient.Get(context.Background(), types.NamespacedName{Name: application.Name, Namespace: application.Namespace}, f)
		}, timeout, interval).ShouldNot(Succeed())
	})
})

var _ = Describe("Application Envs", func() {
	defer GinkgoRecover()

	// generate an application with a single component
	generateApplication := func() *v1alpha1.Application {
		app := generateEmptyApplication()
		app.Spec.Components = append(app.Spec.Components, v1alpha1.ComponentSpec{
			Name:  "test",
			Image: "nginx:latest",
			Env: []v1alpha1.EnvVar{
				{
					Name:  "foo",
					Value: "bar",
					Type:  v1alpha1.EnvVarTypeStatic,
				},
			},
		})
		return app
	}

	Context("Only Static Envs", func() {
		It("", func() {
			By("Create Application")
			application := generateApplication()
			Expect(k8sClient.Create(context.Background(), application)).Should(Succeed())

			deploymentList := &v1.DeploymentList{}
			Eventually(func() bool {
				_ = k8sClient.List(context.Background(), deploymentList, client.MatchingLabels{"kapp-application": application.Name})
				return len(deploymentList.Items) == 1
			}, timeout, interval).Should(Equal(true))

			deployment := deploymentList.Items[0]
			Expect(deployment.Name).Should(Equal(getDeploymentName(application.Name, "test")))
			Expect(*deployment.Spec.Replicas).Should(Equal(int32(1)))
			Expect(len(deployment.Spec.Template.Spec.Containers)).Should(Equal(1))
			Expect(len(deployment.Spec.Template.Spec.Containers[0].Env)).Should(Equal(1))
			Expect(deployment.Spec.Template.Spec.Containers[0].Env[0].Value).Should(Equal("bar"))
		})
	})
})
