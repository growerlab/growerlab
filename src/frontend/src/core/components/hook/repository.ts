import useSWR, { Fetcher } from "swr";
import { useParams } from "react-router-dom";

import { useRepositoryAPI } from "../../api/repository/repository";
import { useGlobal } from "../../global/global";
import {
  RepositoryEntity,
  RepositoryPath,
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

export function useGetRepository(
  pg: RepositoryPathGroup
): Promise<RepositoryEntity> {
  const repositoryAPI = useRepositoryAPI(pg.namespace);

  const fetcher: Fetcher<RepositoryEntity, RepositoryPath> = () => {
    return repositoryAPI.get(pg.repo).then((res) => {
      return res.data.repository;
    });
  };

  const { data, error } = useSWR<RepositoryEntity>(
    `/swr/key/repo/${pg.namespace}/${pg.repo}`,
    fetcher,
    { revalidateOnFocus: false }
  );
  if (error) return Promise.reject(error);
  if (!data)
    return Promise.reject(`pull repo data for ${pg.namespace}/${pg.repo}`);
  return Promise.resolve(data);
}
