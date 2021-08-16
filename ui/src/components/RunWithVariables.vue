<template>
  <v-dialog v-model="dialog" max-width="800px">
    <template v-slot:activator="{ on }">
      <v-btn v-bind="$attrs" v-on="on">Run with Variables</v-btn>
    </template>
    <v-card>
      <v-card-title>Run Workflow with Variables</v-card-title>
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
        <v-divider></v-divider>
        <AddVariableForm :variables.sync="variables" ref="add" />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn text @click="dialog = false">Close</v-btn>
        <v-btn text @click="run">Run Workflow</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import AddVariableForm from "@/components/AddVariableForm";

export default {
  name: "RunWithVariables",

  components: {
    AddVariableForm,
  },

  data: () => ({
    dialog: false,
    types: ["variable", "file"],
    variables: [],
    variable: null,
    value: null,
    masked: false,
    type: "variable",
    attrs: {
      class: "mt-2",
    },
  }),

  props: ["workflow"],

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
    close() {
      this.$refs.add.reset();
      this.dialog = false;
    },
    run() {
      this.$store
        .dispatch("createPipeline", {
          id: this.$props.workflow.id,
          variables: this.variables,
        })
        .then((pipeline) => {
          this.close();
          this.$router.push({
            name: "PipelinesView",
            params: { id: pipeline.id },
          });
        });
    },
  },
};
</script>
