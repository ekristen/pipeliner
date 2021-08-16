<template>
  <v-sheet elevation="2" class="mt-4">
    <v-form class="ma-3">
      <v-text-field
        label="New Workflow Name"
        single-line
        full-width
        v-model="name"
        class="mb-3"
        :error-messages="nameErrors"
        @input="$v.name.$touch()"
        @blur="$v.name.$touch()"
      ></v-text-field>

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
        class="mb-5"
        v-model="workflow"
        :options="cmOptions"
        @ready="onCmReady"
        @focus="onCmFocus"
        @input="onCmCodeChange"
      ></codemirror>

      <div class="pb-4">
        <v-btn
          @click="submit"
          :disabled="submitStatus === 'PENDING' || valid === false"
          class="mr-4"
          color="primary"
        >Create</v-btn>
        <v-btn
          @click="validate"
          class="mr-4"
          color="success"
          :disabled="workflow == ''"
        >Validate/Lint</v-btn>
        <v-btn @click="reset" :disabled="workflow == ''" class="mr-4">clear</v-btn>
        <v-btn @click="loadExample">Load Example Workflow</v-btn>
      </div>
    </v-form>
  </v-sheet>
</template>

<style scoped>
@import "~codemirror/lib/codemirror.css";
@import "~codemirror/addon/lint/lint.css";
@import "~codemirror/theme/monokai.css";
</style>

<script>
import { validationMixin } from "vuelidate";
import { required } from "vuelidate/lib/validators";
import { codemirror } from "vue-codemirror";

import "codemirror/mode/yaml/yaml";
import "codemirror/addon/lint/lint";
import "codemirror/addon/lint/yaml-lint";

window.jsyaml = require("js-yaml"); // Introduce js-yaml to improve core support for grammar checking for codemirror

export default {
  name: "Add",

  components: {
    codemirror,
  },

  mixins: [validationMixin],

  validations: {
    name: {
      required,
      validCharacters(name) {
        return /^[a-z]{1,}[a-z0-9-_]{0,}$/.test(name);
      },
    },
    workflow: { required },
  },

  data() {
    return {
      name: null,
      workflow: "",
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
      },
      submitStatus: null,
    };
  },

  computed: {
    codemirror() {
      return this.$refs.myCm.codemirror;
    },
    nameErrors() {
      const errors = [];
      if (!this.$v.name.$dirty) return errors;
      !this.$v.name.validCharacters && errors.push("name format is invalid");
      !this.$v.name.required && errors.push("name is required.");
      return errors;
    },
    workflowErrors() {
      const errors = [];
      if (!this.$v.workflow.$dirty) return errors;
      !this.$v.workflow.required && errors.push("workflow is required.");
      return errors;
    },
  },

  methods: {
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
      this.workflow = newCode;
    },
    async validate() {
      try {
        await this.$store.dispatch("lintWorkflow", { data: this.workflow });
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
    reset() {
      this.$v.$reset();
      this.valid = false;
      this.errors = [];
      this.name = "";
      this.workflow = "";
      this.submitStatus = null;
    },
    submit() {
      this.$v.$touch();
      if (this.$v.$invalid) {
        this.submitStatus = "ERROR";
      } else {
        this.submitStatus = "PENDING";

        let postData = {
          name: this.name,
          data: this.workflow,
        };

        this.$store
          .dispatch("createWorkflow", { data: postData })
          .then(() => {
            this.submitStatus = "OK";
            this.reset();

            this.$router.push({ path: "/workflows" });
          })
          .catch((err) => {
            this.submitStatus = "ERROR";
            console.log(err);
            this.reset();
          });
      }
    },
    loadExample() {
      this.name = "example-workflow";
      this.workflow = `stages:
  - build

example:
  stage: build
  image: ubuntu:20.04
  script:
    - echo "hello-world"`;
    },
  },
};
</script>
