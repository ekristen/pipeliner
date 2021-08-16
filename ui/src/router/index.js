import Vue from 'vue'
import VueRouter from 'vue-router'
import Dashboard from '../views/dashboard/index.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: {
      breadcrumb: 'Dashboard'
    },
  },
  {
    path: '/workflows',
    name: 'Workflows',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "workflows-list" */ '../views/workflows/list.vue'),
    meta: {
      breadcrumb: 'Workflows'
    },
  },
  {
    path: '/workflows/add',
    name: 'AddWorkflow',
    component: () => import(/* webpackChunkName: "workflows-add" */ '../views/workflows/add.vue'),
    meta: {
      breadcrumb: 'Add'
    }
  },
  {
    path: '/workflows/:id',
    name: 'WorkflowsView',
    component: () => import(/* webpackChunkName: "workflows-view" */ '../views/workflows/view.vue'),
    props: true,
    meta: {
      breadcrumb: 'Workflows'
    }
  },
  {
    path: '/pipelines',
    name: 'Pipelines',
    component: () => import(/* webpackChunkName: "pipelines-list" */ '../views/pipelines/list.vue'),
    meta: {
      breadcrumb: 'Pipelines'
    }
  },
  {
    path: '/pipelines/:id',
    name: 'PipelinesView',
    component: () => import(/* webpackChunkName: "pipelines-view" */ '../views/pipelines/view.vue'),
    props: true,
    meta: {
      breadcrumb: 'View'
    }
  },
  {
    path: '/builds',
    name: 'Builds',
    component: () => import(/* webpackChunkName: "builds-list" */ '../views/builds/list.vue'),
  },
  {
    path: '/builds/:id',
    name: 'BuildsView',
    component: () => import(/* webpackChunkName: "builds-view" */ '../views/builds/view.vue'),
  },
  {
    path: '/builds/:id/artifacts/:artifact_id',
    name: 'ArtifactsView',
    component: () => import(/* webpackChunkName: "artifacts-view" */ '../views/artifacts/view.vue'),
  },
  {
    path: '/runners',
    name: 'runners',
    component: () => import(/* webpackChunkName: "runners-list" */ '../views/runners/list.vue')
  },
  {
    path: '/runners/:id',
    name: 'RunnersView',
    component: () => import(/* webpackChunkName: "runners-view" */ '../views/runners/view.vue')
  },
  {
    path: '/variables',
    name: 'variables',
    component: () => import(/* webpackChunkName: "variables-list" */ '../views/variables/list.vue')
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
