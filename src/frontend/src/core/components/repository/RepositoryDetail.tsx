import React, { Fragment, useMemo, useState } from "react";
import {
  EuiSpacer,
  EuiTab,
  EuiTabbedContentProps,
  EuiTabs,
  EuiText,
} from "@elastic/eui";

import { RepositoryPathGroup, RepositoryEntity } from "../../common/types";
import { Item } from "./Item";
import { useGlobal } from "../../global/global";
import { useRepositoryAPI } from "../../api/repository/repository";
import useSWR, { Fetcher } from "swr";
import Loading from "../common/Loading";
import { Files } from "./detail/Files";
import i18n from "../../i18n/i18n";

export function RepositoryDetail(props: RepositoryPathGroup) {
  const { namespace, repo } = props;
  const global = useGlobal();
  const [currentTab, setCurrentTab] = useState("files");
  const [repository, setRepository] = useState<RepositoryEntity>();

  const repositoryAPI = useRepositoryAPI(namespace);

  const tabs: EuiTabbedContentProps["tabs"] = [
    {
      id: "files",
      name: i18n.t<string>("repository.files"),
      content: <Files branch="master" namespace={namespace} repo={repo} />,
    },
    {
      id: "clone",
      name: i18n.t<string>("repository.clond_and_download"),
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
  const selectedTabContent = useMemo(() => {
    return tabs.find((obj) => obj.id === currentTab)?.content;
  }, [currentTab]);

  const fetcher: Fetcher = () => {
    return repositoryAPI.get(repo).then((res) => {
      setRepository(res.data);
      return res.data;
    });
  };
  useSWR(`/swr/key/repo/${namespace}/${repo}`, fetcher);
  if (!repository) {
    return <Loading />;
  }

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
      <Item global={global} repo={repository} />
      <div className=" mb-3"></div>
      <EuiTabs size="s" className="flex justify-between">
        {renderTabs()}
      </EuiTabs>
      <div className="pt-4 p-1">{selectedTabContent}</div>
    </div>
  );
}
