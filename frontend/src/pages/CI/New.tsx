import { Box, createStyles, Grid, Theme, withStyles, WithStyles } from "@material-ui/core";
import { createDeployAccessTokenAction } from "actions/deployAccessToken";
import { setSuccessNotificationAction } from "actions/notification";
import { push } from "connected-react-router";
import { DeployAccessTokenForm } from "forms/DeployAccessToken";
import { withDeployAccessTokens, WithDeployAccessTokensProps } from "hoc/withDeployAccessTokens";
import { BasePage } from "pages/BasePage";
import React from "react";
import { DeployAccessToken } from "types/deployAccessToken";

const styles = (theme: Theme) =>
  createStyles({
    root: {},
  });

interface Props extends WithStyles<typeof styles>, WithDeployAccessTokensProps {}

interface State {}

class DeployAccessTokenNewPageRaw extends React.PureComponent<Props, State> {
  private submit = async (config: DeployAccessToken) => {
    const { dispatch } = this.props;
    await dispatch(createDeployAccessTokenAction(config));
    this.onSubmitSuccess(config);
    return;
  };

  private onSubmitSuccess = async (config: DeployAccessToken) => {
    const { dispatch } = this.props;
    dispatch(setSuccessNotificationAction("Create Deploy key Successfully"));
    dispatch(push("/ci/keys/" + config.name));
  };

  public render() {
    return (
      <BasePage secondHeaderRight={"New Deploy Key"}>
        <Box p={2}>
          <Grid container spacing={2}>
            <Grid item md={8}>
              <DeployAccessTokenForm onSubmit={this.submit} />
            </Grid>
          </Grid>
        </Box>
      </BasePage>
    );
  }
}

export const DeployAccessTokenNewPage = withStyles(styles)(withDeployAccessTokens(DeployAccessTokenNewPageRaw));
