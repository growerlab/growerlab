// const Router = {
//   Home: {
//     Index: '/',
//     Register: '/register',
//     Login: '/login',
//     ActivateUser: '/activate_user/:code',
//   },
//   User: {
//     Index: '/user/',
//     Repository: {
//       New: '/user/repos/new',
//       List: '/user/repos',
//       Show: '/user/repos/:repoPath',
//     },
//   },
//   Namespace: {
//     Repository: '/:namespacePath/:repoPath',
//   },
// };
//
// module.exports = Router;

// import { IConfig } from 'umi-types';
// import Router from './src/router';
//
// // ref: https://umijs.org/config/
// const config: IConfig = {
//   history: 'hash',
//   treeShaking: true,
//   routes: [
//     {
//       path: Router.Home.Index,
//       component: '../layouts/base',
//       routes: [
//         { path: Router.Home.Index, component: './index' },
//         { path: Router.Home.Register, component: './register' },
//         { path: Router.Home.ActivateUser, component: './activate_user' },
//         { path: Router.Home.Login, component: './login' },
//         // user
//         { path: Router.User.Index, component: './user/index' },
//         { path: Router.User.Repository.New, component: './user/repository/new' },
//         { path: Router.User.Repository.List, component: './user/repository/index' },
//         { path: Router.User.Repository.Show, component: './user/repository/show' },
//         // namespace => repository
//         { path: Router.Namespace.Repository, component: './namespace/repository/show' },
//       ],
//     },
//   ],
//   plugins: [
//     // ref: https://umijs.org/plugin/umi-plugin-react.html
//     [
//       'umi-plugin-react',
//       {
//         antd: true,
//         dva: true,
//         dynamicImport: false,
//         title: 'Rethinking Git',
//         dll: false,
//
//         routes: {
//           exclude: [
//             /models\//,
//             /services\//,
//             /model\.(t|j)sx?$/,
//             /service\.(t|j)sx?$/,
//             /components\//,
//           ],
//         },
//       },
//     ],
//   ],
//   // proxy: {
//   //   '/api': {
//   //     target: 'http://localhost:8080',
//   //     pathRewrite: { '^/api': '' },
//   //     changeOrigin: true,
//   //   },
//   // },
// };
//
// export default config;
