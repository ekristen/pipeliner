<template>
  <v-form>
    <v-text-field
      label="Variable*"
      required
      v-model="variable"
      :error-messages="variableErrors"
      @input="$v.variable.$touch()"
      @blur="$v.variable.$touch()"
    ></v-text-field>
    <v-textarea
      label="Value*"
      required
      v-model="value"
      :error-messages="valueErrors"
      @input="$v.value.$touch()"
      @blur="$v.value.$touch()"
    ></v-textarea>
    <v-select :items="types" label="Type" v-model="type"></v-select>
    <v-checkbox v-model="masked" label="Mask Variable" :disabled="type === 'file'"></v-checkbox>
    <v-btn class="success" @click="addVariable">Add Variable</v-btn>
  </v-form>
</template>

<script>
import { validationMixin } from "vuelidate";
import { required } from "vuelidate/lib/validators";

export default {
  name: "AddVariableForm",

  props: ["submit", "variables"],

  mixins: [validationMixin],

  validations: {
    variable: {
      required,
      validCharacters(variable) {
        return /^[a-zA-Z_]{1,}[a-zA-Z0-9_]{0,}$/.test(variable);
      },
    },
    value: { required },
    type: { required },
  },

  data() {
    return {
      dialog: false,

      variable: "",
      value: "",
      type: "variable",
      masked: false,

      types: ["variable", "file"],

      submitStatus: null,
    };
  },

  computed: {
    variableErrors() {
      const errors = [];
      if (!this.$v.variable.$dirty) return errors;
      !this.$v.variable.validCharacters &&
        errors.push("Variable format is invalid");
      !this.$v.variable.required && errors.push("Variable is required.");

      let existsIndex = this.variables.findIndex(
        (v) => v.name === this.variable
      );
      if (existsIndex !== -1) {
        errors.push("Variable name must be unique");
      }

      return errors;
    },
    valueErrors() {
      const errors = [];
      if (!this.$v.value.$dirty) return errors;
      !this.$v.value.required && errors.push("Variable is required.");
      return errors;
    },
  },

  methods: {
    addVariable() {
      this.$v.$touch();
      if (this.$v.$invalid) {
        return;
      }

      let vName = this.variable;
      let vValue = this.value;
      let vType = this.type;
      let vMasked = this.masked;

      this.variables.push({
        name: vName,
        value: vValue,
        type: vType,
        masked: vMasked,
      });

      this.reset();
    },
    reset() {
      this.$v.$reset();
      this.variable = "";
      this.value = "";
      this.type = "variable";
      this.masked = false;
      this.submitStatus = null;
      this.dialog = false;
    },
    add() {},
  },
};
</script>