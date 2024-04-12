/* Layout */
import Layout from '@/layout'
import store from '@/store'

const Routes = [
  // {
  //   path: 'userinfo',
  //   name: 'userInfo',
  //   component: () => import('@/views/settings/secret'),
  //   meta: { title: '个人中心', icon: 'personal', 'group': 'settings', object: 'cluster' }
  // },
  {
    path: 'secret',
    name: 'settinsSecret',
    component: () => import('@/views/settings/secret'),
    meta: { title: '密钥管理', icon: 'settings_secret', 'group': 'settings', object: 'cluster' }
  },
  {
    path: 'image',
    name: 'settinsImage',
    component: () => import('@/views/settings/image_registry'),
    meta: { title: '镜像仓库', icon: 'docker', 'group': 'settings', object: 'cluster' }
  },
  {
    path: 'spacelet',
    name: 'settinsSpacelet',
    component: () => import('@/views/settings/spacelet'),
    meta: { title: 'Spacelet', icon: 'spacelet', 'group': 'settings', object: 'cluster' }
  },
  // {
  //   path: 'ldap',
  //   name: 'settinsLdap',
  //   component: () => import('@/views/settings/ldap'),
  //   meta: { title: 'Ldap管理', icon: 'ldap', 'group': 'settings', object: 'cluster' }
  // },
  {
    path: 'sonarQube',
    name: 'settinsSonarQube',
    component: () => import('@/views/settings/sonar_qube'),
    meta: { title: 'SonarQube', icon: 'sonar_qube', 'group': 'settings', object: 'cluster' }
  },
  {
    path: 'member',
    name: 'member',
    component: () => import('@/views/settings/member/index'),
    meta: { title: '用户管理', icon: 'member', 'group': 'settings', object: 'user' }
  },
  {
    path: 'platform_role',
    name: 'platform_role',
    component: () => import('@/views/settings/platform_role'),
    meta: { title: '平台权限', icon: 'platform_perm', 'group': 'settings', object: 'role' }
  },
  {
    path: 'audit',
    name: 'platform_audit',
    component: () => import('@/views/settings/audit'),
    meta: { title: '操作审计', icon: 'audit', 'group': 'settings', object: 'role' }
  },
]

const settingsRoutes = [{
  path: 'settings',
  component: Layout,
  hidden: true,
  children: Routes
}]

export default settingsRoutes
