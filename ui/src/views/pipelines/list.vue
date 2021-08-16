<template>
  <v-sheet elevation="2" class="mt-4">
    <v-data-table :headers="headers" :items="pipelines" disable-sort>
      <template v-slot:item.state="{item}">
        <StateChip :state="item.state" :to="{ name: 'PipelinesView', params: { id: item.id } }" />
      </template>
      <template v-slot:item.id="{item}">
        <router-link :to="{ name: 'PipelinesView', params: { id: item.id } }">#{{ item.id }}</router-link>
      </template>
      <template v-slot:item.workflow_id="{item}">
        <router-link
          :to="{name: 'WorkflowsView',params: { id: item.workflow_id }}"
        >{{ getWorkflowName(workflows, item.workflow_id) }}</router-link>
      </template>
      <template v-slot:item.stages="{item}">
        <PipelineStages :id="item.id" />
      </template>
      <template v-slot:item.duration="{item}">{{ item.duration }} seconds</template>
    </v-data-table>
  </v-sheet>
</template>

<script>
import StateChip from "@/components/StateChip";
import PipelineStages from "@/components/PipelineStages";

import { mapState } from "vuex";

export default {
  name: "PipelinesList",

  data: () => ({
    page: 1,
    pageCount: 1,
    headers: [
      {
        text: "Status",
        value: "state",
        width: "1%",
      },
      {
        text: "Pipeline",
        value: "id",
      },
      {
        text: "Workflow",
        value: "workflow_id",
      },
      {
        text: "Stages",
        value: "stages",
      },
      {
        text: "Duration",
        value: "duration",
      },
    ],
  }),

  created() {
    this.$store.dispatch("getPipelines");
    this.$store.dispatch("getWorkflowNames");
  },

  components: {
    StateChip,
    PipelineStages,
  },

  computed: mapState({
    pipelines: (state) => state.pipelines,
    workflows: (state) =>
      state.workflows.map((w) => {
        return { id: w.id, name: w.name };
      }),
  }),

  methods: {
    fgcolor(status) {
      if (status == "success") {
        return "white";
      }

      return "black";
    },
    bgcolor(status) {
      if (status == "success") {
        return "green";
      } else if (status == "skipped") {
        return "grey";
      } else if (status == "running") {
        return "orange";
      } else if (status == "created") {
        return "grey";
      }

      return "red";
    },
    icon(status) {
      if (status == "success") {
        return "mdi-checkbox-marked-circle";
      } else if (status == "skipped") {
        return "mdi-debug-step-over";
      } else if (status == "failed") {
        return "mdi-minus-circle-outline";
      } else if (status == "running") {
        return "mdi-play-circle-outline";
      } else if (status == "pending") {
        return "mdi-dots-horizontal-circle-outline";
      } else if (status == "created") {
        return "mdi-plus-circle-outline";
      }

      return "mdi-help-circle-outline";
    },
    getWorkflowName(workflows, id) {
      return workflows.find((w) => w.id === id).name;
    },
  },
};
</script>
