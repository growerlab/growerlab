import i18n from '../i18n/i18n';

export function getTitle(title?: string): string {
  if (!title || title === "") {
    return i18n.t('website.title');
  }
  title = title + ' | ' + i18n.t('website.title');
  return title
}

