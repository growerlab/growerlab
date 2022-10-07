export const userRules = {
  pwdMinLength: 8,
  pwdMaxLength: 32,
  usernameMinLength: 4,
  usernameMaxLength: 40,
  passwordRegex: /^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$/s,
  usernameRegex: /^[a-zA-Z0-9_-]+$/s,
};

export const repositoryRules = {
  repositoryNameRegex: /^[a-zA-Z0-9_\-.]{2,50}$/,
};
