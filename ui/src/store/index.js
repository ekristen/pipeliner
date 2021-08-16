import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import VueAxios from 'vue-axios'
import * as JSONBig from 'json-bigint'
//import axiosETAGCache from 'axios-etag-cache'


const jsonBig = JSONBig({ storeAsString: true })

Vue.use(Vuex)
Vue.use(VueAxios, axios)

const axiosData = axios.create({
  //transformRequest: data => JSONBig.parse(data),
  transformResponse: (data) => {
    return jsonBig.parse(data)
  },
})

export default new Vuex.Store({
  modules: {

  },
  state: {
    socket: {
      isConnected: false,
      message: '',
      reconnectError: false,
    },

    version: '',

    breadcrumbs: [],
    artifacts: [],
    workflows: [],
    pipelines: [],
    builds: [],
    variables: [],
    runners: [],
    traces: [],
    stages: [],
    stats: {},
    build_tags: [],
    runner_tags: [],
    artifact_files: [],
  },
  mutations: {
    SET_ARTIFACT(state, artifact) {
      const i = state.artifacts.findIndex((n) => n.id === artifact.id)
      if (i !== -1) {
        Vue.set(state.artifacts, i, artifact)
      } else {
        state.artifacts.push(artifact)
      }
    },
    SET_ARTIFACTS(state, artifacts) {
      artifacts.forEach((artifact) => {
        const i = state.artifacts.findIndex((n) => n.id === artifact.id)
        if (i !== -1) {
          Vue.set(state.artifacts, i, artifact)
        } else {
          state.artifacts.push(artifact)
        }
      })
    },


    SET_WORKFLOW(state, workflow) {
      const i = state.workflows.findIndex((n) => n.id === workflow.id)
      if (i !== -1) {
        Vue.set(state.workflows, i, workflow)
      } else {
        state.workflows.push(workflow)
      }
    },
    SET_WORKFLOWS(state, workflows) {
      state.workflows = workflows
    },
    DELETE_WORKFLOW(state, workflow) {
      const i = state.workflows.findIndex((n) => n.name === workflow.name)
      if (i !== -1) {
        Vue.delete(state.workflows, i)
      }
    },

    SET_PIPELINE(state, pipeline) {
      const i = state.pipelines.findIndex((n) => n.id === pipeline.id)
      if (i !== -1) {
        Vue.set(state.pipelines, i, pipeline)
      } else {
        state.pipelines.push(pipeline)
      }
    },
    SET_PIPELINES(state, pipelines) {
      state.pipelines = pipelines
    },

    SET_BUILD(state, build) {
      const i = state.builds.findIndex((n) => n.id === build.id)
      if (i !== -1) {
        Vue.set(state.builds, i, build)
      } else {
        state.builds.push(build)
      }
    },
    SET_BUILDS(state, builds) {
      builds.forEach(build => {
        const i = state.builds.findIndex((n) => n.id === build.id)
        if (i !== -1) {
          Vue.set(state.builds, i, build)
        } else {
          state.builds.push(build)
        }
      })
    },

    SET_BUILD_TAG(state, tag) {
      const i = state.build_tags.findIndex((n) => n.id === tag.id)
      if (i !== -1) {
        Vue.set(state.build_tags, i, tag)
      } else {
        state.build_tags.push(tag)
      }
    },

    SET_BUILD_TAGS(state, tags) {
      tags.forEach(tag => {
        const i = state.build_tags.findIndex((n) => n.id === tag.id)
        if (i !== -1) {
          Vue.set(state.build_tags, i, tag)
        } else {
          state.build_tags.push(tag)
        }
      })
    },

    SET_STAGES(state, stages) {
      stages.forEach(stage => {
        const i = state.stages.findIndex((e) => ((e.id === stage.id) || (e.pipeline_id === stage.pipeline_id && e.index === stage.index)))
        if (i !== -1) {
          Vue.set(state.stages, i, stage)
        } else {
          state.stages.push(stage)
        }
      });
    },

    SET_STAGE(state, stage) {
      const i = state.stages.findIndex((e) => ((e.id === stage.id) || (e.pipeline_id === stage.pipeline_id && e.index === stage.index)))
      if (i !== -1) {
        Vue.set(state.stages, i, stage)
      } else {
        state.stages.push(stage)
      }
    },

    SET_RUNNERS(state, runners) {
      runners.forEach(runner => {
        const i = state.runners.findIndex((e) => e.id === runner.id)
        if (i !== -1) {
          Vue.set(state.runners, i, runner)
        } else {
          state.runners.push(runner)
        }
      });
    },
    SET_RUNNER(state, runner) {
      const i = state.runners.findIndex((n) => n.id === runner.id)
      if (i !== -1) {
        Vue.set(state.runners, i, runner)
      } else {
        state.runners.push(runner)
      }
    },

    SET_TRACE(state, trace) {
      const i = state.traces.findIndex((n) => n.id === trace.id)
      if (i !== -1) {
        Vue.set(state.traces, i, trace)
      } else {
        state.traces.push(trace)
      }
    },

    SET_VARIABLES(state, variables) {
      state.variables = variables
    },

    SET_VARIABLE(state, variable) {
      const i = state.variables.findIndex((n) => n.name === variable.name)
      if (i !== -1) {
        Vue.set(state.variables, i, variable)
      } else {
        state.variables.push(variable)
      }
    },
    DELETE_VARIABLE(state, variable) {
      const i = state.variables.findIndex((n) => n.name === variable.name)
      if (i !== -1) {
        Vue.delete(state.variables, i)
      }
    },

    SET_STATS(state, stats) {
      state.stats = stats;
    },

    SET_RUNNER_TAG(state, tag) {
      const i = state.runner_tags.findIndex((n) => n.id === tag.id)
      if (i !== -1) {
        Vue.set(state.runner_tags, i, tag)
      } else {
        state.runner_tags.push(tag)
      }
    },

    SET_RUNNER_TAGS(state, tags) {
      tags.forEach(tag => {
        const i = state.runner_tags.findIndex((n) => n.id === tag.id)
        if (i !== -1) {
          Vue.set(state.runner_tags, i, tag)
        } else {
          state.runner_tags.push(tag)
        }
      })
    },

    SET_ARTIFACT_CONTENTS(state, files) {
      files.forEach(file => {
        const i = state.artifact_files.findIndex((f) => f.artifact_id == file.artifact_id && f.id == file.id)
        if (i !== -1) {
          Vue.set(state.artifact_files, i, file)
        } else {
          state.artifact_files.push(file)
        }
      })
    },

    SET_VERSION(state, version) {
      state.version = version.version;
    },

    RESET_BREADCRUMB(state) {
      state.breadcrumbs = []
    },
    ADD_BREADCRUMB(state, { text, href }) {
      state.breadcrumbs.push({
        text, href,
      })
    },

    SOCKET_ONOPEN(state, event) {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true

      /*
      mutations.DELETE_ALERT(state, {
        id: 'websocket-disconnected',
      })
      */
    },
    SOCKET_ONCLOSE(state, event) {
      state.socket.isConnected = false
      console.log(event)
      /*
      mutations.SET_ALERT(state, {
        id: 'websocket-disconnected',
        color: 'warning',
        description: 'The websocket connection to the backend has been lost.',
      })
      */
    },
    SOCKET_ONERROR(state, event) {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE(state, message) {
      state.socket.message = message
      console.log('message received: ', message)
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT(state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR(state) {
      state.socket.reconnectError = true;
    },
  },
  actions: {
    addBuild({ commit }, { data }) {
      commit('SET_BUILD', data)
    },
    addPipeline({ commit }, { data }) {
      commit('SET_PIPELINE', data)
    },
    addStage({ commit }, { data }) {
      commit('SET_STAGE', data)
    },
    addRunner({ commit }, { data }) {
      commit('SET_RUNNER', data)
    },
    addRunnerTag({ commit }, { data }) {
      commit('SET_RUNNER_TAG', data)
    },
    addWorkflow() {

    },
    addArtifact({ commit }, { data }) {
      commit('SET_ARTIFACT', data)
    },

    addBuildTag({ commit }, { data }) {
      commit('SET_BUILD_TAG', data)
    },

    addGlobalVariable({ commit }, { data }) {
      commit('SET_VARIABLE', data)
    },
    deleteGlobalVariable({ commit }, { data }) {
      commit('DELETE_VARIABLE', data)
    },

    createGlobalVariable({ commit }, { data }) {
      return axiosData.post(`${window.location.origin}/v1/variables`, data)
        .then(r => r.data)
        .then(variable => {
          commit('SET_VARIABLE', variable)
        })
    },
    removeGlobalVariable({ commit }, { name }) {
      return axiosData.delete(`${window.location.origin}/v1/variables/${name}`)
        .then(r => r.data)
        .then(variable => {
          commit('DELETE_VARIABLE', variable)
        })
    },

    getWorkflows({ commit }) {
      axiosData.get(`${window.location.origin}/v1/workflows`)
        .then(r => r.data)
        .then(workflows => {
          commit('SET_WORKFLOWS', workflows)
        })
    },

    getWorkflowNames({ commit }) {
      axiosData.get(`${window.location.origin}/v1/workflows?fields=name`)
        .then(r => r.data)
        .then(workflows => {
          commit('SET_WORKFLOWS', workflows)
        })
    },

    getWorkflow({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/workflows/${id}`)
        .then(r => r.data)
        .then(workflow => {
          commit('SET_WORKFLOW', workflow)
        })
    },
    createWorkflow({ commit }, { data }) {
      axiosData.post(`${window.location.origin}/v1/workflows`, data)
        .then(r => r.data)
        .then(workflow => {
          commit('SET_WORKFLOW', workflow)
        })
    },
    removeWorkflow({ commit }, { id }) {
      axiosData.delete(`${window.location.origin}/v1/workflows/${id}`)
        .then(r => r.data)
        .then(workflow => {
          commit('DELETE_WORKFLOW', workflow)
        })
    },
    async updateWorkflow({ commit }, { id, data }) {
      return await axiosData.patch(`${window.location.origin}/v1/workflows/${id}`, { data })
        .then(r => r.data)
        .then(workflow => {
          commit('SET_WORKFLOW', workflow)
        })
    },
    async lintWorkflow(_, { data }) {
      return await axiosData.post(`${window.location.origin}/v1/lint`, {
        data,
      }).then(r => r.data)
    },

    createPipeline({ commit }, { id, variables }) {
      return axiosData.post(`${window.location.origin}/v1/workflows/${id}/pipeline`, {
        variables,
      })
        .then(r => r.data)
        .then(pipeline => {
          commit('SET_PIPELINE', pipeline)
          return pipeline
        })
    },

    getPipeline({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/pipelines/${id}`)
        .then(r => r.data)
        .then(pipeline => {
          commit('SET_PIPELINE', pipeline)
        })
    },

    getPipelines({ commit }) {
      axiosData.get(`${window.location.origin}/v1/pipelines`)
        .then(r => r.data)
        .then(pipelines => {
          commit('SET_PIPELINES', pipelines)
        })
    },

    getPipelineBuilds({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/pipelines/${id}/builds`)
        .then(r => r.data)
        .then(builds => {
          commit('SET_BUILDS', builds)
        })
    },

    getPipelineStages({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/pipelines/${id}/stages`)
        .then(r => r.data)
        .then(stages => {
          commit('SET_STAGES', stages)
        })
    },

    getPipelineArtifacts({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/pipelines/${id}/artifacts`)
        .then(r => r.data)
        .then(artifacts => {
          commit('SET_ARTIFACTS', artifacts)
        })
    },

    getPipelineStageBuilds({ commit }, { pipeline_id, stage_id }) {
      axiosData.get(`${window.location.origin}/v1/pipelines/${pipeline_id}/stages/${stage_id}/builds`)
        .then(r => r.data)
        .then(builds => {
          commit('SET_BUILDS', builds)
        })
    },

    deletePipeline({ commit }, { id }) {
      axiosData.delete(`${window.location.origin}/v1/pipelines/${id}`)
        .then(r => r.data)
        .then(pipeline => {
          commit('REMOVE_PIPELINE', pipeline)
        })
    },
    cancelPipeline({ commit }, { id }) {
      axiosData.post(`${window.location.origin}/v1/pipelines/${id}/cancel`)
        .then(r => r.data)
        .then(pipeline => {
          commit('SET_PIPELINE', pipeline)
        })
    },

    runJob({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/builds/${id}/run`)
        .then(r => r.data)
        .then(build => {
          commit('SET_BUILD', build)
        })
    },

    getBuild({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/builds/${id}`)
        .then(r => r.data)
        .then(build => {
          commit('SET_BUILD', build)
        })
    },

    getBuildArtifact({ commit }, { id, type }) {
      axiosData.get(`${window.location.origin}/v1/builds/${id}/artifact/${type}`)
        .then(r => r.data)
        .then(artifact => {
          commit('SET_ARTIFACT', artifact)
        })
    },

    getBuildTrace({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/builds/${id}/trace`)
        .then(r => r.data)
        .then(trace => {
          commit('SET_TRACE', trace)
        })
    },

    getBuilds({ commit }) {
      return axiosData.get(`${window.location.origin}/v1/builds`)
        .then(r => r.data)
        .then(builds => {
          commit('SET_BUILDS', builds)
        })
    },

    getJobTags({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/builds/${id}/tags`)
        .then(r => r.data)
        .then(tags => {
          commit('SET_BUILD_TAGS', tags)
        })
    },

    getRunner({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/runners/${id}`)
        .then(r => r.data)
        .then(runner => {
          commit('SET_RUNNER', runner)
        })
    },

    getRunners({ commit }) {
      axiosData.get(`${window.location.origin}/v1/runners`)
        .then(r => r.data)
        .then(runners => {
          commit('SET_RUNNERS', runners)
        })
    },

    getVariablesGlobal({ commit }) {
      axiosData.get(`${window.location.origin}/v1/variables`)
        .then(r => r.data)
        .then(variables => {
          commit('SET_VARIABLES', variables)
        })
    },

    getStats({ commit }) {
      axiosData.get(`${window.location.origin}/v1/stats`)
        .then(r => r.data)
        .then(stats => {
          commit('SET_STATS', stats)
        })
    },

    getArtifactFiles({ commit }, { id }) {
      axiosData.get(`${window.location.origin}/v1/artifacts/${id}/files`)
        .then(r => r.data)
        .then(files => {
          console.log(files)
          commit('SET_ARTIFACT_CONTENTS', files)
        })
    },

    async retryJob({ commit }, { id }) {
      return await axiosData.post(`${window.location.origin}/v1/builds/${id}/retry`)
        .then(r => r.data)
        .then(job => {
          commit('SET_BUILD', job)
          return job
        })
    },
    cancelJob({ commit }, { id }) {
      axiosData.post(`${window.location.origin}/v1/builds/${id}/cancel`)
        .then(r => r.data)
        .then(job => {
          commit('SET_BUILD', job)
        })
    },

    getVersion({ commit }) {
      axiosData.get(`${window.location.origin}/v1/version`)
        .then(r => r.data)
        .then(version => {
          commit('SET_VERSION', version)
        })
    },

    saveRunner({ commit }, { id, description, tags, run_untagged }) {
      axiosData.patch(`${window.location.origin}/v1/runners/${id}`, {
        description,
        tags,
        run_untagged,
      })
        .then(r => r.data)
        .then(runner => {
          commit('SET_RUNNER', runner)
        })
    }
  },
  getters: {
    pipelines: (state) => state.pipelines,
    pipelineByID: (state) => (id) => {
      return state.pipelines.find((d) => d.id === id);
    },
    buildsByPipelineIDAndStage: (state) => (pipeline_id, index) =>
      state.builds.filter((b) => b.pipeline_id == pipeline_id && b.stage_idx == index),

    buildsByPipelineID: (state) => (pipeline_id) =>
      state.builds.filter((b) => b.pipeline_id == pipeline_id),
  }
})
