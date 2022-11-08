import React, { useEffect } from "react";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiTextArea,
  EuiSwitch,
  EuiButton,
} from "@elastic/eui";
import { useForm, useWatch } from "react-hook-form";

import i18n from "../../i18n/i18n";
import { repositoryRules } from "../../api/rule";

interface IFormNewRepository {
  namespace_path: string;
  name: string;
  description: string;
  public: boolean;
}

interface IProps {
  ownerPath: string;
}

export function NewRepositoryForm(props: IProps) {
  const {
    register,
    handleSubmit,
    setValue,
    setError,
    control,
    formState: { errors, isValid },
  } = useForm<IFormNewRepository>({
    defaultValues: {
      namespace_path: props.ownerPath,
      public: false,
    },
  });

  const fields = useWatch({ control });

  register("name", {
    required: i18n.t<string>("notice.required"),
    pattern: {
      value: repositoryRules.repositoryNameRegex,
      message: i18n.t<string>("repository.tooltip.name"),
    },
    validate: (value) => {
      return value.indexOf("..") !== 0;
    },
  });
  register("description", {
    required: false,
    maxLength: 255,
  });
  register("public");

  useEffect(() => {
    console.info(fields);
  }, [fields]);

  const onSubmit = (data: IFormNewRepository) => {
    console.log(data);
    //       Message.Success(t("repository.tooltip.success"));
    //       router.push(Router.User.Repository.List);
  };

  return (
    <EuiForm component="form" onSubmit={handleSubmit(onSubmit)}>
      <EuiFormRow
        label={i18n.t<string>("repository.name")}
        isInvalid={!!errors.name}
        error={errors.name?.message}
      >
        <EuiFieldText
          type={"text"}
          icon={"editorCodeBlock"}
          isInvalid={false}
          onChange={(e) => {
            setValue("name", e.target.value, { shouldValidate: true });
            if (e.target.value.indexOf("..") === 0) {
              setError("name", {
                type: "onChange",
                message: i18n.t<string>("repository.tooltip.name"),
              });
            }
          }}
          placeholder={i18n.t<string>("repository.name")}
        />
      </EuiFormRow>

      <EuiFormRow
        label={i18n.t<string>("repository.description")}
        isInvalid={!!errors.description}
        error={errors.description?.message}
      >
        <EuiTextArea
          required={false}
          isInvalid={false}
          onChange={(e) => {
            setValue("description", e.target.value, { shouldValidate: true });
          }}
        />
      </EuiFormRow>

      <EuiFormRow label={i18n.t<string>("repository.public")}>
        <EuiSwitch
          name={"public"}
          label=""
          checked={true}
          onChange={(e) => {
            setValue("public", !e.target.checked);
            return;
          }}
        />
      </EuiFormRow>

      <EuiButton type="submit" disabled={!isValid}>
        {i18n.t<string>("repository.create_repository")}
      </EuiButton>
    </EuiForm>
  );
}
