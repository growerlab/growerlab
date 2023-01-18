import React, { Fragment, Suspense, useMemo, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import {
  EuiSpacer,
  EuiTab,
  EuiTabbedContentProps,
  EuiTabs,
  EuiText,
} from "@elastic/eui";

import { RepositoryPathGroup, RepositoryEntity } from "../../common/types";
import { Header } from "./Header";
import { useGlobal } from "../../global/global";
import { useRepositoryAPI } from "../../api/repository";
import useSWR from "swr";
import Loading from "../common/Loading";
import { Files } from "./detail/Files";
import i18n from "../../i18n/i18n";

export function RepositoryDetail(props: RepositoryPathGroup) {
  const { namespace, repo } = props;
  const global = useGlobal();
  const [currentTab, setCurrentTab] = useState("files");

  const [, setSearchParams] = useSearchParams();
  const { folder } = useParams();

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
          reference="main"
          folder={folder || ""}
          namespace={namespace}
          repo={repo}
          repository={data}
          onChangeReference={(reference: string) => {
            setSearchParams({ ref: reference });
          }}
          onChangeFilePath={(filePath: string) => {
            if (filePath !== "" && filePath !== "/") {
              setSearchParams({ filePath: filePath });
            } else {
              setSearchParams({ filePath: "" });
            }
          }}
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
