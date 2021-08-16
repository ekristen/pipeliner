<template>
  <v-sheet elevation="2" class="mt-4 pa-4" v-if="workflow">
    <span class="text-h4">{{ workflow.name }}</span>

    <div class="float-right">
      <Run :workflow="workflow" class="mr-2" color="primary" />
      <RunWithVariables :workflow="workflow" class="mr-2" color="secondary" />
    </div>
    <v-divider></v-divider>

    <v-alert dense text type="success" v-if="valid">This workflow is valid.</v-alert>

    <v-alert
      border="top"
      color="red lighten-2"
      dark
      v-for="(data, index) in errors"
      :key="index"
      v-html="data"
    ></v-alert>

    <codemirror
      class="mb-5 mt-5"
      v-model="code"
      :options="cmOptions"
      @ready="onCmReady"
      @focus="onCmFocus"
      @input="onCmCodeChange"
    ></codemirror>

    <div class="mt-2">
      <v-btn v-show="edit == false" :disabled="edit" @click="edit = !edit" class="mr-2">Edit</v-btn>
      <v-btn
        @click="validate"
        v-show="edit == true && valid == false"
        class="mr-2"
        color="success"
        :disabled="edit == false"
      >Validate/Lint</v-btn>
      <v-btn
        v-show="edit == true && valid == true"
        :disabled="!edit"
        @click="save"
        class="mr-2"
        color="primary"
      >Save</v-btn>
      <v-btn v-show="edit == true" :disabled="!edit" @click="cancel" color="warning">Cancel</v-btn>
    </div>
  </v-sheet>
</template>

<style scoped>
@import "~codemirror/lib/codemirror.css";
@import "~codemirror/addon/lint/lint.css";
@import "~codemirror/theme/monokai.css";
</style>

<style>
.CodeMirror {
  height: 100% !important;
}
</style>

<script>
import { mapState } from "vuex";

import Run from "@/components/Run";
import RunWithVariables from "@/components/RunWithVariables";

import { codemirror } from "vue-codemirror";
import "codemirror/mode/yaml/yaml";
import "codemirror/addon/lint/lint";
import "codemirror/addon/lint/yaml-lint";

window.jsyaml = require("js-yaml");

export default {
  name: "WorkflowsView",

  data() {
    return {
      edit: false,
      code: "",
      valid: false,
      errors: [],
      cmOptions: {
        lineNumbers: true, // display line number
        mode: "text/x-yaml", // grammar model
        gutters: ["CodeMirror-lint-markers"], // Syntax checker
        theme: "monokai", // Editor theme
        lint: true, // Turn on grammar checking
        tabSize: 2,
        indentUnit: 2,
        smartIndent: true,
        insertSoftTab: true,
        readOnly: "nocursor",
      },
      originalWorkflow: "",
    };
  },

  components: {
    codemirror,
    Run,
    RunWithVariables,
  },

  created() {
    this.fetchData();
  },

  watch: {
    $route: "fetchData",
    workflow(value) {
      this.code = value.data;
      this.originalWorkflow = value.data;
    },
    edit(value) {
      if (value) {
        this.cmOptions.readOnly = false;
      } else {
        this.cmOptions.readOnly = "nocursor";
      }
    },
  },

  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror;
    },
    ...mapState({
      workflow(state) {
        return state.workflows.find((d) => d.id === this.$route.params.id);
      },
    }),
  },

  methods: {
    async save() {
      try {
        await this.$store.dispatch("updateWorkflow", {
          id: this.$route.params.id,
          data: this.code,
        });

        this.cancel();
      } catch (err) {
        console.log(err);
      }
    },
    async validate() {
      try {
        await this.$store.dispatch("lintWorkflow", { data: this.code });
        this.valid = true;
        this.errors = [];
      } catch (err) {
        this.valid = false;
        const data = err.response.data;
        if (typeof data.errors != "undefined") {
          this.errors = data.errors;
        }
      }
    },
    cancel() {
      this.edit = false;
      this.valid = false;
      this.code = this.originalWorkflow;
    },
    onCmReady(cm) {
      cm.setOption("extraKeys", {
        Tab: function (cm) {
          var spaces = Array(cm.getOption("indentUnit") + 1).join(" ");
          cm.replaceSelection(spaces);
        },
      });
    },
    onCmFocus() {},
    onCmCodeChange(newCode) {
      this.code = newCode;
    },
    fetchData() {
      const workflowId = this.$route.params.id;
      this.$store.dispatch("getWorkflow", { id: workflowId });
    },
  },
};
</script>
