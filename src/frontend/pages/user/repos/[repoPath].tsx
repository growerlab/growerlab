import React, {useEffect, useState} from 'react';
// @ts-ignore
import {FormComponentProps} from '@ant-design/compatible/lib/form';
import {Empty, Menu, PageHeader, Popover, Tag} from 'antd';

import {RepositoryDetail} from '../../../components/repository/RepositoryDetail';
import {getTitle} from '../../../common/document';
import {Repository} from '../../../api/repository/repository';
import {Message} from '../../../api/common/notice';
import i18n from '../../../i18n';

interface RepoPath {
  repoPath: string;
}

export default function (props: FormComponentProps) {
  const {repoPath} = props.match.params as RepoPath;

  getTitle(i18n.t(repoPath));

  return (
    <div>
      <RepositoryDetail repoPath={repoPath}/>
    </div>
  );
}
