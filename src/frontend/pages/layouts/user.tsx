import React, { useState, useEffect, useRef } from "react";
import Link from "next/link";
import { useRouter } from "next/router";
import { withTranslation } from "react-i18next";
import { HomeIcon, ProjectsIcon, GitRepoIcon, IconComponent, BuggyIcon, SearchInput, Popover, Menu, TrashIcon, EditIcon, CircleArrowRightIcon, PeopleIcon, Button } from 'evergreen-ui'

import { Session } from "../../services/auth/session";
import { Message } from "../../api/common/notice";
import { Router } from "../../config/router";




function UserLayout(props: any) {
  const { t } = props;
  const router = useRouter();

  useEffect((): void => {
    // 验证用户是否登录
    Session.isLogin().catch(() => {
      Message.Warning(t("user.tooltip.not_login"));
      router.push(Router.Home.Login);
    });
  });

  // eslint-disable-next-line react-hooks/rules-of-hooks
  const [collapsed, setCollapsed] = useState(false);
  const plusMenu = (
    <Link href={Router.User.Repository.New}>{t("repository.new")}</Link>
  );

  const logoutClick = (): void => {
    Session.logout(router);
  };

  const userMenu = (
    <div>
      <Popover
        position={'bottom-left'}
        content={
          <Menu>
            <Menu.Group title="Actions">
              {[
                ["Your Profile", "/", PeopleIcon],
                ["Settings", "/", CircleArrowRightIcon],
              ].map(([title, url, icon]) => (
                <Menu.Item icon={icon as IconComponent}>
                  <a href={url as string}>
                    {title}
                  </a>
                </Menu.Item>
              ))}
            </Menu.Group>
            <Menu.Divider />
            <Menu.Group title="destructive">
              <Menu.Item icon={TrashIcon} intent="danger">
                <Link passHref href={""}>
                  <a
                    onClick={() => logoutClick()}
                  >
                    {t("user.logout")}
                  </a>
                </Link>
              </Menu.Item>
            </Menu.Group>
          </Menu>
        }
      >
        <img
          className="h-8 w-8 rounded-full"
          src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
          alt=""
        />
      </Popover>
    </div>
  );

  // const path = window.location.pathname.split('/').slice(0, 3);
  // const menuKey = [path.join('/')];
  const MenuItem = (props: { icon: IconComponent, title: string, href: string, selected?: boolean }) => {
    const iconRef = useRef(null);

    return (
      <div>
        <div className="opacity-70 hover:opacity-100">
          <a href={props.href} className="text-white block px-4 py-3 rounded-md text-sm hover:bg-blue-900 text-center">
            <div className=" mb-2">
              {props.icon}
            </div>
            {props.title}
          </a>
        </div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex flex-row fixed bottom-0 w-full top-0">
        <div className="bg-blue-800 ">
          <nav className="flex flex-col h-full">
            <div className="flex-none">
              <div>
                <a
                  href="#"
                  className="text-white block px-3 py-5  text-base font-medium text-center"
                  aria-current="page"
                >
                  <BuggyIcon size={45} className="inline" />
                </a>
              </div>
              <div className="px-2 pt-2 pb-3 space-y-1 ">
                {[
                  ["Home", Router.User.Index, <HomeIcon size={20} className="inline" />],
                  [t("repository.list"), Router.User.Repository.List, <GitRepoIcon size={20} className="inline" />],
                  ["Projects", "/", <ProjectsIcon size={20} className="inline" />],
                ].map(([title, href, icon]) => (
                  // <menuItem icon={icon} title={title} href={href}></menuItem>
                  <MenuItem icon={icon} href={href} title={title}></MenuItem>
                ))}
              </div>
            </div>
            <div className="flex-auto">{/* 填充 */}</div>

          </nav>
        </div>

        <div className="grow">
          <div className="flex flex-col h-full">
            <header className="bg-white shadow">
              <div className="max-w-full  mx-auto py-3 px-4 sm:px-2 lg:px-6">
                <div className="flex columns-2">
                  <div className="flex-none">
                    <SearchInput placeholder="Search..." />
                  </div>
                  <div className="grow"></div>
                  <div className="flex-none">
                    {userMenu}
                  </div>
                </div>
              </div>
            </header>
            <main>
              <div className="max-w-full mx-auto py-4 sm:px-4 lg:px-8">
                <div className="px-4 py-6 sm:px-0">
                  <div className="border-0 border-dashed border-gray-200 rounded-lg max-h-full">
                    {props.children}
                  </div>
                </div>
              </div>
            </main>
          </div>
        </div>
      </div>
    </div>
  );
}

export default withTranslation()(UserLayout);
