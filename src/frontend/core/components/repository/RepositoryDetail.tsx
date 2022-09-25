import React, { useEffect, useState } from "react";

import { LoginInfo } from "../../services/auth/session";
import { Repository } from "../../services/repository/repository";
import { RepositoryArgs } from "../../services/repository/types";
import { repoIcon } from "./common";

interface RepositoryDetailProps extends RepositoryArgs {
  currentUser?: LoginInfo; // 当前登录的用户，不一定是仓库的所有者
}

export function RepositoryDetail(props: RepositoryDetailProps) {
  const { ownerPath, repoPath } = props;
  const [current, setCurrent] = useState("code");

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
