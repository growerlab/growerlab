import React from "react";
import { EuiLoadingContent } from "@elastic/eui";
import { useTitle } from "react-use";

import { getTitle } from "../../common/document";

export default function Loading() {
  useTitle(getTitle("Loading..."));

  return (
    <div>
      <EuiLoadingContent lines={2} />
    </div>
  );
}
