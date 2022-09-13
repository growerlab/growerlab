import React, { useState } from "react";
import {
  SettingOutlined,
  CodeOutlined,
  IssuesCloseOutlined,
  CloudDownloadOutlined,
} from "@ant-design/icons";
import { Menu, PageHeader, Popover, Tag, Tabs, Input, Empty } from "antd";
import { LockOutlined, UnlockOutlined } from "@ant-design/icons/lib";

import { RepositoryArgs } from "../../api/repository/types";
import { Repository } from "../../api/repository/repository";
import { repoIcon } from "./common";

export function RepositoryDetail(props: RepositoryArgs) {
  const { repoPath } = props;
  const [current, setCurrent] = useState("code");
  const { TabPane } = Tabs;

  const repo = new Repository({ repoPath: repoPath });
  const repoData = repo.get();
  if (repoData === null) {
    return (
      <div>
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
      </div>
    );
  }
  const repository = repoData.repository;

  const handleClick = (e: any) => {
    console.log(e.key);
    if (e.key === "clone") {
      return;
    }
    setCurrent(e.key);
  };

  return (
    <div>
      <PageHeader
        title={
          <span>
            {repoIcon(repoData.repository.public)}
            {repository.pathGroup}
          </span>
        }
      />
      <Menu onClick={handleClick} selectedKeys={[current]} mode="horizontal">
        <Menu.Item key="code">
          <CodeOutlined />
          Code
        </Menu.Item>
        <Menu.Item key="issues" disabled={true}>
          <IssuesCloseOutlined />
          Issues
        </Menu.Item>
        <Menu.Item key="settings" style={{ float: "right" }} disabled={true}>
          <SettingOutlined />
          Settings
        </Menu.Item>
        <Menu.Item key="clone" style={{ float: "right" }} onBlur={() => {}}>
          <CloudDownloadOutlined />
          <Popover
            placement="bottom"
            title="Clone or download"
            content={
              <Tabs defaultActiveKey="1" size={"small"}>
                <TabPane tab="Http" key="1" active={true}>
                  <Input
                    placeholder="Basic usage"
                    defaultValue={repository.gitHttpURL}
                    readOnly={true}
                  />
                  {/* <span>{repository.gitHttpURL}</span> */}
                </TabPane>
                <TabPane tab="SSH" key="2" animated={false}>
                  <Input
                    placeholder="Basic usage"
                    defaultValue={repository.gitSshURL}
                    readOnly={true}
                  />
                  {/* <span>{repository.gitSshURL}</span> */}
                </TabPane>
              </Tabs>
            }
          >
            <span>Clone or download</span>
          </Popover>
        </Menu.Item>
      </Menu>
    </div>
  );
}
