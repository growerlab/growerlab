import React, { useState } from "react";
import i18n from "i18next";
import { EuiCodeBlock, EuiCard, EuiCallOut, EuiIcon } from "@elastic/eui";

import CloneURL from "./CloneURL";
import { useTitle } from "../../../global/state";

interface Props {
  defaultBranch: string;
  cloneURLSSH: string;
  cloneURLHttp: string;
}

export default function EmptyTree(props: Props) {
  useTitle(i18n.t("repository.empty_tree"));

  const [cloneURL, setCloneURL] = useState<string>(props.cloneURLHttp);

  const newShell = `echo "# hello" >> README.md 
git init
git add README.md
git commit -m 'first commit'
git remote add origin ${cloneURL || props.cloneURLHttp}
git push -u origin ${props.defaultBranch}`;

  const oldShell = `git remote add origin ${cloneURL || props.cloneURLHttp}
git push -u origin ${props.defaultBranch}`;

  return (
    <div>
      <EuiCard
        title={""}
        display="plain"
        hasBorder
        description={
          <div className={"text-left"}>
            <h2>
              <EuiIcon type={"tokenRepo"} size={"xl"} />
              {i18n.t<string>("repository.empty_tree")}
            </h2>

            <EuiCallOut
              title={i18n.t<string>("repository.empty_good_news")}
              iconType="pin"
              className={"mb-6 mt-6"}
            >
              <p>{i18n.t<string>("repository.empty_good_news_detail")}</p>
            </EuiCallOut>

            <CloneURL {...props} onChange={(url) => setCloneURL(url)} />

            <h4>
              <strong>Option 1: </strong>
              {i18n.t<string>("repository.empty_option1")}
            </h4>
            <EuiCodeBlock language="bash" fontSize="s" paddingSize="m">
              {newShell}
            </EuiCodeBlock>
            <h4>
              <strong>Option 2: </strong>
              {i18n.t<string>("repository.empty_option2")}
            </h4>
            <EuiCodeBlock language="bash" fontSize="s" paddingSize="m">
              {oldShell}
            </EuiCodeBlock>
          </div>
        }
      />
    </div>
  );
}
