<template>
  <div>
    <v-simple-table>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">Type</th>
            <th class="text-left">Name</th>
            <th class="text-left">Value</th>
            <th class="text-left">Masked</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in variables" :key="item.name">
            <td>{{ item.file == true ? "file" : "variable" }}</td>
            <td>{{ item.name }}</td>
            <td>{{ item.value }}</td>
            <td>{{ item.masked }}</td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "Variables",

  props: ["scope", "id"],

  updated() {
    console.log(this.scope);
    if (this.scope == "workflow") {
      // this
    } else if (this.scope == "pipeline") {
      // this
    } else {
      this.$store.dispatch("getVariablesGlobal");
    }
  },

  computed: mapState({
    variables: (state) => state.variables,
  }),
};
</script>
