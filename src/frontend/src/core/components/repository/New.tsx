import React from "react";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiTextArea,
  EuiSwitch,
  EuiButton,
} from "@elastic/eui";

import { repositoryRules } from "../../api/rule";
import i18n from "../../i18n/i18n";

interface NewRepositoryPayload {
  namespacePath: string;
  name: string;
  public: boolean;
}

export function NewRepositoryFrom(props: any) {
  const { t } = props;

  // const [createRepository, { loading: mutationLoading, error: mutationError }] =
  //   useMutation<{
  //     input: NewRepositoryPayload;
  //   }>(GQL_REGISTER, {
  //     onCompleted: (data: any) => {
  //       Message.Success(t("repository.tooltip.success"));
  //       router.push(Router.User.Repository.List);
  //     },
  //   });

  // const handleSubmit = (values: Store) => {
  //   const payload = values as NewRepositoryPayload;
  //
  //   // createRepository({
  //   //   variables: {
  //   //     input: payload,
  //   //   },
  //   // });
  // };

  return (
    <EuiForm component="form">
      <EuiFormRow
        label={i18n.t<string>("repository.owner")}
        // isInvalid={emailValidateMsg != null}
        // error={emailValidateMsg}
      >
        <EuiFieldText
          name={"namespace_path"}
          type={"text"}
          required={true}
          onChange={(event) => {
            // setEmail(event.target.value);
          }}
          icon={"user"}
          isInvalid={false}
          // onBlur={onBlur}
        />
      </EuiFormRow>

      <EuiFormRow
        label={i18n.t<string>("repository.name")}
        // isInvalid={emailValidateMsg != null}
        // error={emailValidateMsg}
      >
        <EuiFieldText
          name={"name"}
          type={"text"}
          required={true}
          onChange={(event) => {
            // setEmail(event.target.value);
          }}
          icon={"user"}
          isInvalid={false}
          // onBlur={onBlur}
        />
      </EuiFormRow>

      <EuiFormRow
        label={i18n.t<string>("repository.description")}
        // isInvalid={emailValidateMsg != null}
        // error={emailValidateMsg}
      >
        <EuiTextArea
          name={"description"}
          required={true}
          onChange={(event) => {
            // setEmail(event.target.value);
          }}
          isInvalid={false}
          // onBlur={onBlur}
        />
      </EuiFormRow>

      <EuiFormRow
        label={i18n.t<string>("repository.public")}
        // isInvalid={emailValidateMsg != null}
        // error={emailValidateMsg}
      >
        <EuiSwitch
          name={"public"}
          label=""
          checked={false}
          onChange={(e) => {
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
