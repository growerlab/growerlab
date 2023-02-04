import React, { useState } from "react";
import {
  copyToClipboard,
  EuiButtonGroup,
  EuiButtonIcon,
  EuiFieldText,
  EuiFlexGroup,
} from "@elastic/eui";

import { useNotice } from "../../../global/state/useNotice";
import i18n from "../../../i18n/i18n";

interface Props {
  defaultBranch: string;
  cloneURLSSH: string;
  cloneURLHttp: string;
  onChange(url: string): void;
}

export default function CloneURL(props: Props) {
  const [cloneURL, setCloneURL] = useState<string>(props.cloneURLHttp);
  const notice = useNotice();

  const toggleButtons = [
    {
      id: `http`,
      label: "HTTP",
      value: props.cloneURLHttp,
    },
    {
      id: `ssh`,
      label: "SSH",
      value: props.cloneURLSSH,
    },
  ];

  const [toggleIdSelected, setToggleIdSelected] = useState("http");

  const onChange = (optionId: string, optionValue: string) => {
    setToggleIdSelected(optionId);
    setCloneURL(optionValue);
    props.onChange(optionValue);
  };

  return (
    <>
      <EuiFlexGroup>
        <EuiButtonGroup
          legend="This is a basic group"
          options={toggleButtons}
          idSelected={toggleIdSelected}
          onChange={(id, value) => onChange(id, value)}
        />
        <EuiFieldText
          value={cloneURL}
          readOnly={true}
          compressed={true}
          append={
            <EuiButtonIcon
              iconType={"copy"}
              onClick={() => {
                copyToClipboard(cloneURL);
                notice.success(i18n.t<string>("notice.copy_notice"));
              }}
            />
          }
        />
      </EuiFlexGroup>
    </>
  );
}
