<template>
  <v-sheet elevation="2" class="mt-4">
    <v-btn :to="{ name: 'AddWorkflow' }" class="ma-2">Add Workflow</v-btn>
    <v-simple-table>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">Name</th>
            <th class="text-left">Updated At</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="workflow in workflows" :key="workflow.id">
            <td>
              <router-link
                :to="{ name: 'WorkflowsView', params: { id: workflow.id } }"
              >{{ workflow.name }}</router-link>
            </td>
            <td>{{ workflow.updated_at }}</td>
            <td>
              <Run :workflow="workflow" class="mr-2" color="primary" small />
              <RunWithVariables :workflow="workflow" class="mr-2" color="secondary" small />
              <v-btn color="error" class="ml-2" small @click="deleteWorkflow(workflow.id)">Delete</v-btn>
            </td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </v-sheet>
</template>

<script>
import Run from "@/components/Run";
import RunWithVariables from "@/components/RunWithVariables";

import { mapState } from "vuex";

export default {
  name: "Workflows",

  components: {
    Run,
    RunWithVariables,
  },

  created() {
    this.$store.dispatch("getWorkflows");
  },

  computed: mapState({
    workflows: (state) => state.workflows,
  }),

  methods: {
    deleteWorkflow(id) {
      this.$store.dispatch("removeWorkflow", { id });
    },
  },
};
</script>
