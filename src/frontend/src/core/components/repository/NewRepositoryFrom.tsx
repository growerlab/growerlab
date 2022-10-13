import React, { useEffect, useState } from "react";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiTextArea,
  EuiSwitch,
  EuiButton,
} from "@elastic/eui";

import i18n from "../../i18n/i18n";
import { useForm, useWatch } from "react-hook-form";
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

export function NewRepositoryFrom(props: IProps) {
  const [privateChecked, setPrivateChecked] = useState(false);

  const {
    register,
    handleSubmit,
    setValue,
    control,
    formState: { errors },
  } = useForm<IFormNewRepository>({
    defaultValues: {
      namespace_path: props.ownerPath,
    },
  });

  const fields = useWatch({ control });

  register("name", {
    required: i18n.t<string>("notice.required"),
    pattern: {
      value: repositoryRules.repositoryNameRegex,
      message: i18n.t<string>("repository.tooltip.name"),
    },
  });
  register("description", {
    required: false,
    maxLength: 255,
  });

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
        isInvalid={errors.name != null}
        error={errors.name?.message}
      >
        <EuiFieldText
          type={"text"}
          icon={"editorCodeBlock"}
          isInvalid={false}
          onChange={(e) => {
            setValue("name", e.target.value, { shouldValidate: true });
          }}
        />
      </EuiFormRow>

      <EuiFormRow
        label={i18n.t<string>("repository.description")}
        isInvalid={errors.description != null}
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
          checked={privateChecked}
          onChange={(e) => {
            setPrivateChecked(!privateChecked);
            setValue("public", !privateChecked);
            return;
          }}
        />
      </EuiFormRow>

      <EuiButton type="submit" fill>
        {i18n.t<string>("repository.create_repository")}
      </EuiButton>
    </EuiForm>
  );
}
