import React, { Fragment, useState } from "react";
import {
  EuiSpacer,
  EuiTab,
  EuiTabbedContentProps,
  EuiTabs,
  EuiText,
} from "@elastic/eui";

import { RepositoryPathGroup, RepositoryEntity } from "../../common/types";
import { repoIcon } from "./common";
import { useGetRepository } from "../hook/repository";

export function RepositoryDetail(props: RepositoryPathGroup) {
  const { namespace, repo } = props;
  const [currentTab, setCurrentTab] = useState("code");
  const [repository, setRepository] = useState<RepositoryEntity>();

  const repoEntity = useGetRepository({ namespace: namespace, repo: repo });
  repoEntity.then((data) => {
    console.info(data);
    setRepository(data);
  });

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

      <EuiTabs size="s" className="flex justify-between">
        {renderTabs()}
      </EuiTabs>
    </div>
  );
}
