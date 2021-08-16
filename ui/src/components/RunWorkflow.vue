<template>
  <v-dialog v-model="dialog" max-width="800px">
    <template v-slot:activator="{ on, attrs }">
      <v-btn color="primary" v-bind="attrs" v-on="on" class="mr-2" small>{{ title }}</v-btn>
    </template>
    <v-card>
      <v-card-title>
        <span class="headline">Run Workflow: {{ workflow.name }}</span>
      </v-card-title>
      <v-card-text>
        <v-simple-table v-if="variables.length > 0">
          <template v-slot:default>
            <thead>
              <tr>
                <th class="text-left">Name</th>
                <th class="text-left">Value</th>
                <th class="text-left">Type</th>
                <th class="text-left">Masked</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in variables" :key="item.name">
                <td>{{ item.name }}</td>
                <td>{{ item.value }}</td>
                <td>{{ item.type }}</td>
                <td>{{ item.masked }}</td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
        <v-form v-if="this.withVariables">
          <v-text-field label="Variable*" required v-model="variable"></v-text-field>
          <v-textarea label="Value*" required v-model="value"></v-textarea>
          <v-select :items="types" label="Type" v-model="type"></v-select>
          <v-checkbox v-model="masked" label="Mask Variable"></v-checkbox>
          <v-btn color="success" class="mr-4" @click="addVariable">Add Variable</v-btn>
          <small>*indicates required field</small>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="dialog = false">Close</v-btn>
        <v-btn text @click="run">Run</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import router from "@/router";

export default {
  name: "RunWorkflow",

  data: () => ({
    dialog: false,
    types: ["variable", "file"],
    variables: [],
    variable: null,
    value: null,
    masked: false,
    type: "variable",
  }),

  props: ["workflow", "withVariables", "title"],

  updated() {
    if (this.dialog == true && this.withVariables == "false") {
      this.run();
    }
  },

  methods: {
    addVariable() {
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

      this.variable = "";
      this.value = "";
      this.type = "variable";
      this.masked = false;
    },
    run() {
      this.$store
        .dispatch("createPipeline", { id: this.$props.workflow.id })
        .then((pipeline) => {
          router.push({
            name: "PipelinesView",
            params: { id: pipeline.id },
          });
          this.dialog = false;
        });
    },
  },
};
</script>
