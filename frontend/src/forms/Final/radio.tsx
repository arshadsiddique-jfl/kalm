import React from "react";
import { FieldRenderProps } from "react-final-form";
import { FormLabel } from "@material-ui/core";
import FormControl from "@material-ui/core/FormControl";
import FormControlLabel from "@material-ui/core/FormControlLabel";
import Radio from "@material-ui/core/Radio";
import RadioGroup from "@material-ui/core/RadioGroup";
import { Body2 } from "widgets/Label";

interface KRadioGroupRenderOption {
  value: string;
  label: React.ReactNode;
  explain?: string;
}

interface FinalRadioGroupRenderProps extends FieldRenderProps<string> {
  options: KRadioGroupRenderOption[];
  title?: string;
}

export const FinalRadioGroupRender = ({
  title,
  options,
  input: { onChange, value },
  meta: { error },
}: FinalRadioGroupRenderProps) => {
  return (
    <FormControl component="fieldset" fullWidth margin="dense" error={error}>
      {title ? <FormLabel component="legend">{title}</FormLabel> : null}
      <RadioGroup
        value={value}
        onChange={(_: any, value: string) => {
          onChange(value);
        }}
      >
        {options.map((option) => {
          return (
            <span key={option.value}>
              <FormControlLabel value={option.value} control={<Radio />} label={option.label} />
              {option.explain ? <Body2 style={{ padding: "0 16px 0 32px" }}>{option.explain}</Body2> : null}
            </span>
          );
        })}
      </RadioGroup>
    </FormControl>
  );
};
