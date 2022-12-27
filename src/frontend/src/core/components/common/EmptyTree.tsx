import React from "react";
import { useTitle } from "react-use";
import { EuiEmptyPrompt, EuiButton } from "@elastic/eui";

import { getTitle } from "../../common/document";

export default function EmptyTree(props: any) {
  useTitle(getTitle("Empty repository..."));

  return (
    <div>
      <EuiEmptyPrompt
        iconType="securityAnalyticsApp"
        iconColor="default"
        title={<h2>Start adding cases</h2>}
        titleSize="xs"
        body={<p>Add a new case or change your filter settings.</p>}
        actions={
          <EuiButton size="s" color="primary" fill>
            Add a case
          </EuiButton>
        }
      />
    </div>
  );
}
