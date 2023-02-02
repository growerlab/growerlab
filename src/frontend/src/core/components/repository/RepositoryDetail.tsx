import React, { Fragment, Suspense, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import {
  EuiSpacer,
  EuiTab,
  EuiTabbedContentProps,
  EuiTabs,
  EuiText,
} from "@elastic/eui";
import useSWR from "swr";

import {
  RepositoryPathGroup,
  RepositoryEntity,
  DetailType,
} from "../../common/types";
import { Header } from "./Header";
import { useGlobal } from "../../global/global";
import { useRepositoryAPI } from "../../api/repository";
import Loading from "../ui/common/Loading";
import { Files } from "./files/Files";
import i18n from "../../i18n/i18n";

export type Props = RepositoryPathGroup;

export function RepositoryDetail(props: Props) {
  const { namespace, repo } = props;
  const global = useGlobal();
  const [currentTab, setCurrentTab] = useState("files");
  const path = useParams()["*"];
  const refType = useParams()["refType"] as DetailType;

  const repositoryAPI = useRepositoryAPI(namespace);

  const fetcher = () =>
    repositoryAPI.getDetail(repo).then((res) => {
      return res.data;
    });
  const { data } = useSWR<RepositoryEntity>(
    `/swr/key/repo/${namespace}/${repo}`,
    fetcher,
    { suspense: true }
  );

  const tabs: EuiTabbedContentProps["tabs"] = [
    {
      id: "files",
      name: i18n.t<string>("repository.files"),
      content: (
        <Files
          type={refType}
          blobPath={path}
          reference="main"
          initialFolder={path || ""}
          namespace={namespace}
          repo={repo}
          repository={data}
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
  }, [currentTab, namespace, repo]);

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
      <Suspense fallback={<Loading lines={3} />}>
        <Header global={global} repo={data!} />
        <div className=" mb-3"></div>
        <EuiTabs size="s" className="flex justify-between">
          {renderTabs()}
        </EuiTabs>
        <div className="pt-4">{selectedTabContent}</div>
      </Suspense>
    </div>
  );
}
