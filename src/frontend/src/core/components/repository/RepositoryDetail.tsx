import {
  EuiSpacer,
  EuiTab,
  EuiTabbedContentProps,
  EuiTabs,
  EuiText,
} from "@elastic/eui";
import React, { Fragment, useEffect, useState } from "react";
import useSWR, { Fetcher } from "swr";

import { UserInfo } from "../../services/auth/session";
import { Repository } from "../../services/repository/repository";
import {
  RepositoryArgs,
  RepositoryEntity,
} from "../../services/repository/types";
import { repoIcon } from "./common";

interface RepositoryDetailProps extends RepositoryArgs {
  currentUser?: UserInfo; // 当前登录的用户，不一定是仓库的所有者
}

const fetcher: Fetcher<RepositoryEntity, RepositoryArgs> = (
  args: RepositoryArgs
) => {
  const repo = new Repository(args.ownerPath);
  return repo.get(args.repoPath).then((res) => {
    return res.data.repository;
  });
};

export function RepositoryDetail(props: RepositoryDetailProps) {
  const { ownerPath, repoPath } = props;
  const [currentTab, setCurrentTab] = useState("code");
  const [repository, setRepository] = useState<RepositoryEntity>();

  const { data } = useSWR<RepositoryEntity>(
    { ownerPath: ownerPath, repoPath: repoPath },
    fetcher,
    { suspense: true }
  );
  setRepository(data);

  // const repo = new Repository(ownerPath);
  // repo.get(repoPath).then((res) => {
  //   setRepository(res.data.repository);
  // });

  const tabs: EuiTabbedContentProps["tabs"] = [
    {
      id: "code",
      name: "Code",
      content: (
        <Fragment>
          <EuiSpacer />
          <EuiText>
            <p>
              Cobalt is a chemical element with symbol Co and atomic number 27.
              Like nickel, cobalt is found in the Earth&rsquo;s crust only in
              chemically combined form, save for small deposits found in alloys
              of natural meteoric iron. The free element, produced by reductive
              smelting, is a hard, lustrous, silver-gray metal.
            </p>
          </EuiText>
        </Fragment>
      ),
    },
    {
      id: "clone",
      name: "Clone or download",
      content: (
        <Fragment>
          <EuiSpacer />
          <EuiText>
            <p>
              Intravenous sugar solution, also known as dextrose solution, is a
              mixture of dextrose (glucose) and water. It is used to treat low
              blood sugar or water loss without electrolyte loss.
            </p>
          </EuiText>
        </Fragment>
      ),
    },
  ];
  const renderTabs = () => {
    return tabs.map((tab, index) => (
      <EuiTab
        key={index}
        onClick={() => setCurrentTab(tab.id)}
        isSelected={tab.id === currentTab}
        disabled={tab.disabled}
        prepend={tab.prepend}
        append={tab.append}
      >
        {tab.name}
      </EuiTab>
    ));
  };
  return (
    <div>
      <h3 className={"text-xl"}>
        {repoIcon(repository?.public ? true : false)}
        <span className="ml-4"></span>
        {repository?.path}
      </h3>
      <div className={"mt-6 mb-4 text-gray-400"}>{repository?.description}</div>

      <EuiTabs>{renderTabs()}</EuiTabs>
    </div>
  );
}
