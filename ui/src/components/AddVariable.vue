<template>
  <v-row>
    <v-dialog v-model="dialog" persistent max-width="800px">
      <template v-slot:activator="{ on, attrs }">
        <v-btn color="primary" class="ml-3" dark v-bind="attrs" v-on="on">Add Variable</v-btn>
      </template>

      <v-card>
        <v-card-title>
          <span class="headline">Add Global Variable</span>
        </v-card-title>
        <v-card-text>
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
            <small>*indicates required field</small>
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="warning darken-1" text @click="dialog = false">Close</v-btn>
          <v-btn
            color="primary darken-1"
            text
            @click="submit"
            :disabled="submitStatus === 'PENDING'"
          >Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>

<script>
import { validationMixin } from "vuelidate";
import { required } from "vuelidate/lib/validators";

export default {
  name: "AddVariable",

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
    reset() {
      this.$v.$reset();
      this.variable = "";
      this.value = "";
      this.type = "variable";
      this.masked = false;
      this.submitStatus = null;
      this.dialog = false;
    },
    submit() {
      this.$v.$touch();
      if (this.$v.$invalid) {
        this.submitStatus = "ERROR";
      } else {
        this.submitStatus = "PENDING";

        let postData = {
          name: this.variable,
          value: this.value,
          masked: this.masked,
          file: false,
        };

        if (this.type === "file") {
          postData.file = true;
        }

        this.$store
          .dispatch("createGlobalVariable", { data: postData })
          .then(() => {
            this.submitStatus = "OK";
            this.reset();
          })
          .catch((err) => {
            this.submitStatus = "ERROR";
            console.log(err);
            this.reset();
          });
      }
    },
  },
};
</script>