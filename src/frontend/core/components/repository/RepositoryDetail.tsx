import {
  EuiSpacer,
  EuiTab,
  EuiTabbedContentProps,
  EuiTabs,
  EuiText,
} from "@elastic/eui";
import React, { Fragment, useState } from "react";

import { LoginInfo } from "../../services/auth/session";
import { Repository } from "../../services/repository/repository";
import { RepositoryArgs } from "../../services/repository/types";
import { repoIcon } from "./common";

interface RepositoryDetailProps extends RepositoryArgs {
  currentUser?: LoginInfo; // 当前登录的用户，不一定是仓库的所有者
}

export function RepositoryDetail(props: RepositoryDetailProps) {
  const { ownerPath, repoPath } = props;
  const [currentTab, setCurrentTab] = useState<"code" | "clone">("code");

  const repo = new Repository(ownerPath);
  const repoData = repo.get(repoPath);
  if (repoData === null) {
    return (
      <div>
        <h3 className="text-xl">404 - Not found</h3>
      </div>
    );
  }
  const repository = repoData.repository;

  const handleClick = (e: any) => {
    console.log(e.key);
    if (e.key === "clone") {
      return;
    }
    setCurrentTab(e.key);
  };

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
        href={tab.href}
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
        {repoIcon(repoData.repository.public)}
        <span className="ml-4"></span>
        {repository.path}
      </h3>
      <div className={"mt-6 mb-4 text-gray-400"}>{repository.description}</div>

      <EuiTabs>{renderTabs()}</EuiTabs>
    </div>
  );
}
