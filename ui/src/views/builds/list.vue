<template>
  <v-sheet elevation="2" class="mt-4">
    <v-data-table :headers="headers" :items="builds" disable-sort>
      <template v-slot:item.state="{item}">
        <StateChip :state="item.state" :to="{ name: 'BuildsView', params: { id: item.id } }" />
      </template>
      <template v-slot:item.id="{item}">
        <router-link :to="{ name: 'BuildsView', params: { id: item.id } }">#{{ item.id }}</router-link>
      </template>
      <template v-slot:item.pipeline_id="{item}">
        <router-link
          :to="{name: 'PipelinesView', params: { id: item.pipeline_id }}"
        >#{{ item.pipeline_id }}</router-link>
      </template>
    </v-data-table>
  </v-sheet>
</template>

<script>
import StateChip from "@/components/StateChip";

import { mapState } from "vuex";

export default {
  name: "Builds",

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
        text: "Job",
        value: "id",
      },
      {
        text: "Pipeline",
        value: "pipeline_id",
      },
      {
        text: "Stage",
        value: "stage",
      },
      {
        text: "Name",
        value: "name",
      },
      {
        text: "Duration",
        value: "duration",
      },
    ],
  }),

  created() {
    this.$store.dispatch("getBuilds");
  },

  components: {
    StateChip,
  },

  computed: mapState({
    builds: (state) =>
      state.builds
        .sort((a, b) => {
          if (a.id < b.id) return -1;
          if (a.id > b.id) return 1;
          return 0;
        })
        .reverse(),
  }),

  methods: {
    id(entity) {
      if (entity.id === entity.id_str) {
        return entity.id;
      } else {
        return entity.id_str;
      }
    },

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
      }

      return "pink";
    },
    icon(status) {
      if (status == "success") {
        return "mdi-checkbox-marked-circle";
      } else if (status == "skipped") {
        return "mdi-debug-step-over";
      } else if (status == "failed") {
        return "mdi-minus-circle-outline";
      }

      return "mdi-help-circle-outline";
    },
  },
};
</script>
