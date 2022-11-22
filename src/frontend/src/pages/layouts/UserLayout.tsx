import React, { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";

import {
  EuiFieldSearch,
  EuiIcon,
  EuiPopover,
  EuiContextMenu,
  useGeneratedHtmlId,
  EuiContextMenuPanelDescriptor,
  EuiButtonIcon,
} from "@elastic/eui";

import { Router } from "../../config/router";
import { Session } from "../../core/services/auth/session";
import { useGlobal } from "../../core/global/init";
import i18n from "../../core/i18n/i18n";
import { useTitle } from "react-use";

interface Props extends React.PropsWithChildren {
  title: string;
}

export default function UserLayout(props: Props) {
  const { title } = props;
  const { notice } = useGlobal();
  const navigate = useNavigate();

  useTitle(title);

  useEffect((): void => {
    // 验证用户是否登录
    Session.isLogin().catch(() => {
      notice!.warning(i18n.t("user.tooltip.not_login"));
      navigate(Router.Home.Login);
    });
  }, []);

  const [collapsed, setCollapsed] = useState(false);
  const plusMenu = (
    <EuiButtonIcon
      display="base"
      iconType="plus"
      color={"primary"}
      size={"s"}
      title={i18n.t<string>("repository.new")}
      onClick={() => navigate(Router.User.Repository.New)}
    />
  );

  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const onOpenUserMenu = () => setIsPopoverOpen(!isPopoverOpen);
  const onCloseUserMenu = () => setIsPopoverOpen(false);
  const userMenuButton = (
    <img
      className="h-8 w-8 rounded-full"
      src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
      alt=""
      onClick={onOpenUserMenu}
    />
  );

  const contextMenuPopoverId = useGeneratedHtmlId({
    prefix: "contextMenuPopover",
  });

  const panels: EuiContextMenuPanelDescriptor[] = [
    {
      id: 0,
      title: "This is a context menu",
      items: [
        {
          name: "Your Profile",
          icon: <EuiIcon type="user" />,
          href: "/",
        },
        {
          name: "Settings",
          icon: "indexSettings",
          href: "/",
        },
        {
          isSeparator: true,
          key: "sep",
        },
        {
          name: i18n.t<string>("user.logout"),
          onClick: () => Session.logout(),
        },
      ],
    },
  ];

  const userMenu = (
    <>
      <EuiPopover
        id={contextMenuPopoverId}
        button={userMenuButton}
        isOpen={isPopoverOpen}
        closePopover={onCloseUserMenu}
        panelPaddingSize="none"
        anchorPosition="downLeft"
      >
        <EuiContextMenu initialPanelId={0} panels={panels} />
      </EuiPopover>
    </>
  );

  const MenuItem = (props: {
    icon: React.ReactNode;
    title: string;
    href: string;
    selected?: boolean;
  }) => {
    return (
      <div>
        <div className="opacity-70 hover:opacity-100">
          <Link
            to={props.href}
            className="text-white block px-4 py-3 rounded-md text-sm hover:bg-blue-900 text-center"
          >
            <div className=" mb-2">{props.icon}</div>
            {props.title}
          </Link>
        </div>
      </div>
    );
  };

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
                  <EuiIcon type={"color"} className="inline" />
                </a>
              </div>
              <div className="px-2 pt-2 pb-3 space-y-1 ">
                {[
                  [
                    "Home",
                    Router.User.Index,
                    <EuiIcon type="grid" key={"home"} />,
                  ],
                  [
                    i18n.t<string>("repository.menu"),
                    Router.User.Repository.Index,
                    <EuiIcon
                      type={"visVega"}
                      key={"repository"}
                      className="inline"
                    />,
                  ],
                  [
                    i18n.t<string>("project.menu"),
                    Router.User.Project.Index,
                    <EuiIcon
                      type={"sessionViewer"}
                      key={"project"}
                      className="inline"
                    />,
                  ],
                  [
                    i18n.t<string>("thought.menu"),
                    Router.User.Project.Index,
                    <EuiIcon
                      type={"sessionViewer"}
                      key={"project"}
                      className="inline"
                    />,
                  ],
                ].map(([title, href, icon]) => (
                  <MenuItem
                    key={title.toString()}
                    icon={icon}
                    href={href.toString()}
                    title={title.toString()}
                  ></MenuItem>
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
                <div className="flex">
                  <div className="flex-none"></div>
                  {/* search */}
                  <div className="grow">
                    <div className="flex">
                      <div className="flex-1"></div>
                      <div className="flex-1">
                        <EuiFieldSearch
                          placeholder="Search cmd+k"
                          // value={value}
                          // onChange={(e) => onChange(e)}
                          // isClearable={isClearable}
                          fullWidth={true}
                          aria-label="Use aria labels when no actual label is in use"
                        />
                      </div>
                      <div className="flex-1"></div>
                    </div>
                  </div>
                  {/* user */}
                  <div className="flex-none">
                    <div className="flex">
                      <div className={"mr-5"}>{plusMenu}</div>
                      <div>{userMenu}</div>
                    </div>
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
