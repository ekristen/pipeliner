<template>
  <v-sheet elevation="2" class="mt-4">
    <v-data-table
      :headers="headers"
      :items="runners"
      disable-sort
      disable-pagination
      disable-filtering
      hide-default-footer
    >
      <template v-slot:item="{ item }">
        <StateChip :state="state(item)" />
      </template>
      <template v-slot:item="{ item }"
        >{{ item.platform }} / {{ item.architecture }}</template
      >
      <template v-slot:item="{ item }">
        <v-chip-group>
          <Tags :tags="item.tags" color="tag" small />
          <v-chip color="secondary" small label v-if="item.run_untagged"
            >Run Untagged</v-chip
          >
        </v-chip-group>
      </template>
      <template v-slot:item="{ item }">{{
        formatDate(item.contacted_at)
      }}</template>
      <template v-slot:item="{ item }">
        <RunnerEdit :runner="item" />
        <v-btn fab x-small color="danger" class="ml-2">
          <v-icon>mdi-cross</v-icon>
        </v-btn>
      </template>
    </v-data-table>

    <DocsRunners :runners="runners" class="pa-4" />
  </v-sheet>
</template>

<script>
import { mapState } from "vuex";

import RunnerEdit from "@/components/RunnerEdit";
import StateChip from "@/components/StateChip";
import DocsRunners from "@/components/DocsRunners";
import Tags from "@/components/Tags";

import State from "@/mixins/state";

export default {
  name: "RunnersList",

  mixins: [State],

  components: {
    DocsRunners,
    RunnerEdit,
    StateChip,
    Tags,
  },

  created() {
    this.$store.dispatch("getRunners");
  },

  computed: {
    ...mapState({
      runners: (state) =>
        state.runners.sort((a, b) => {
          if (a.contacted_at < b.contacted_at) return 1;
          if (a.contacted_at > b.contacted_at) return -1;

          return 0;
        }),
    }),
  },

  data() {
    return {
      headers: [
        {
          text: "Status",
          align: "start",
          sortable: false,
          value: "status",
        },
        { text: "Name", value: "name" },
        { text: "Description", value: "description" },
        { text: "Platform", value: "platform" },
        { text: "Tags", value: "tags" },
        { text: "Last Contacted", value: "contacted_at" },
        { text: "", value: "actions" },
      ],
    };
  },

  methods: {
    state(runner) {
      if (!runner.active) {
        return "inactive";
      }

      if (runner.contacted_at === null) {
        return "offline";
      }

      const contactedAt = Date.parse(runner.contacted_at);
      const now = new Date();
      const diff = Math.floor((now - contactedAt) / 1000);

      if (diff > 60 * 2) {
        return "offline";
      }

      return "active";
    },
    formatDate(date) {
      if (date === null) {
        return "never";
      }
      return this.$moment(date).fromNow();
    },
  },
};
</script>
