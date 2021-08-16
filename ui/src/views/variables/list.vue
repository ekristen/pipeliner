<template>
  <v-sheet elevation="2" class="mt-4 pt-1">
    <v-alert border="top" colored-border type="info" elevation="2" class="ma-4">
      These are
      <strong>GLOBAL VARIABLES</strong>, these will be automatically
      added to EVERY job of EVERY pipeline, but they can be overwritten by the
      pipeline or job specific variables.
    </v-alert>
    <AddVariable class="ma-1" />
    <v-simple-table>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">Type</th>
            <th class="text-left">Name</th>
            <th class="text-left">Value</th>
            <th class="text-left">Masked</th>
            <th class="text-left">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in variables" :key="item.name">
            <td>{{ item.file == true ? "file" : "variable" }}</td>
            <td>{{ item.name }}</td>
            <td>**********</td>
            <td>{{ item.masked }}</td>
            <td>
              <v-btn color="error" small @click="deleteVariable(item.name)">Delete</v-btn>
            </td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </v-sheet>
</template>

<script>
import { mapState } from "vuex";

import AddVariable from "@/components/AddVariable";

export default {
  name: "Variables",

  created() {
    this.$store.dispatch("getVariablesGlobal");
  },

  components: {
    AddVariable,
  },

  computed: mapState({
    variables: (state) => state.variables,
  }),

  methods: {
    deleteVariable(name) {
      this.$store.dispatch("removeGlobalVariable", { name });
    },
  },
};
</script>
