import { useParams } from "react-router-dom";

import { useGlobal } from "../../global/global";
import {
  RepositoryPathGroup,
} from "../../common/types";

interface PathGroupMaybe extends RepositoryPathGroup {
  isInvalid: boolean;
}

export function useRepositoryPathGroup(): PathGroupMaybe {
  const global = useGlobal();
  const currentUser = global.currentUser;
  const { namespace, repo } = useParams();

  if (
    (currentUser === undefined && namespace === undefined) ||
    repo === undefined
  )
    return { namespace: "", repo: "", isInvalid: true };

  // 优先使用url中的namespace
  if (namespace !== undefined)
    return { namespace: namespace, repo: repo, isInvalid: false };

  if (currentUser === undefined)
    return { namespace: "", repo: "", isInvalid: true };

  return { namespace: currentUser.namespace, repo: repo, isInvalid: false };
}
