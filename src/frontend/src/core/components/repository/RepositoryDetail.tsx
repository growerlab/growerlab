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
import { useRepositoryAPI } from "../../api/repository";
import useSWRImmutable, { Fetcher } from "swr";
import Loading from "../common/Loading";
import { Files } from "./detail/Files";
import i18n from "../../i18n/i18n";

export function RepositoryDetail(props: RepositoryPathGroup) {
  const { namespace, repo } = props;
  const global = useGlobal();
  const [currentTab, setCurrentTab] = useState("files");
  const [repository, setRepository] = useState<RepositoryEntity>();
  const [isEmptyTree, setTreeEmpty] = useState<boolean>(false);

  const repositoryAPI = useRepositoryAPI(namespace);

  const tabs: EuiTabbedContentProps["tabs"] = [
    {
      id: "files",
      name: i18n.t<string>("repository.files"),
      content: (
        <Files
          reference="main"
          dir="/"
          namespace={namespace}
          repo={repo}
          isEmptyTree={isEmptyTree}
        />
      ),
    },
    {
      id: "clone",
      name: i18n.t<string>("repository.clone_and_download"),
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
  }, [currentTab, isEmptyTree]);

  const fetcher: Fetcher = () => {
    return repositoryAPI.getDetail(repo).then((res) => {
      setRepository(res.data);
      return res.data;
    });
  };
  useSWRImmutable(`/swr/key/repo/${namespace}/${repo}`, fetcher);
  if (!repository) {
    return <Loading />;
  }
  if (!isEmptyTree && repository.last_push_at == 0) {
    setTreeEmpty(true);
  }

  const renderTabs = () => {
    return tabs.map((tab, index) => (
      <EuiTab
        key={index}
        onClick={() => setCurrentTab(tab.id)}
        isSelected={tab.id === currentTab}
        disabled={tab.disabled}
        // prepend={tab.prepend}
        // append={tab.append}
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
