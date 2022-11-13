import React, { useEffect, useState } from "react";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiTextArea,
  EuiSwitch,
  EuiButton,
} from "@elastic/eui";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";

import i18n from "../../i18n/i18n";
import { repositoryRules } from "../../api/rule";
import {
  Repository,
  RepositoryRequest,
} from "../../services/repository/repository";
import { useGlobal } from "../../global/init";
import { Router } from "../../../config/router";
import { Namespace } from "../../common/types";

export function NewRepositoryForm(props: Namespace) {
  const [isPublic, setPublic] = useState(true);
  const navigate = useNavigate();
  const { notice } = useGlobal();
  const {
    register,
    handleSubmit,
    setValue,
    setError,
    formState: { errors, isValid },
  } = useForm<RepositoryRequest>({
    defaultValues: {
      namespace: props.namespace,
      public: true,
      description: "",
    },
  });

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

  const onSubmit = (data: RepositoryRequest) => {
    console.log(data);

    const service = new Repository(props.namespace);
    service.create(data).then((res) => {
      notice?.success(i18n.t("repository.tooltip.success"));
      console.info(res);
      navigate(Router.User.Repository.Show.render({ repo: data.name }));
    });
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
          isInvalid={false}
          onChange={(e) => {
            setValue("name", e.target.value, { shouldValidate: true });
            if (e.target.value.indexOf("..") === 0) {
              setError("name", {
                type: "value",
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
          checked={isPublic}
          onChange={(e) => {
            setPublic(e.target.checked);
            setValue("public", e.target.checked);
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
