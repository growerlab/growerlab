import React from "react";
import { EuiControlBar } from "@elastic/eui";

import { EuiBreadcrumbProps } from "@elastic/eui/src/components/breadcrumbs/breadcrumb";

interface Props {
  truncate?: boolean;
  max?: number;
  icon?: string;
  breadcrumbs: EuiBreadcrumbProps[];
}

export default function Breadcrumbs(props: Props) {
  return (
    <div>
      <EuiControlBar
        position={"relative"}
        showContent={false}
        className={"!bg-inherit !shadow-none"}
        controls={[
          {
            iconType: "submodule",
            id: "root_icon",
            controlType: "icon",
            "aria-label": "Project Root",
          },
          {
            controlType: "breadcrumbs",
            id: "current_file_path",
            responsive: true,
            breadcrumbs: props.breadcrumbs,
            truncate: props.truncate || false,
            max: props.max || null,
          },
          // {
          //   controlType: "spacer",
          // },
          // {
          //   controlType: "icon",
          //   id: "branch_icon",
          //   iconType: "branch",
          //   "aria-label": "Branch Icon",
          // },
          // {
          //   controlType: "button",
          //   id: "open_history_view",
          //   label: "切换分支",
          //   color: "primary",
          //   onClick: undefined,
          // },
          // {
          //   controlType: "divider",
          // },
          // {
          //   controlType: "button",
          //   id: "open_history_view",
          //   label: "Show history",
          //   color: "primary",
          //   onClick: undefined,
          // },
        ]}
      />
    </div>
  );
}
