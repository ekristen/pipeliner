<template>
  <v-simple-table>
    <template v-slot:default>
      <thead>
        <tr>
          <th class="text-left">Status</th>
          <th class="text-left">Name</th>
          <th class="text-left">Description</th>
          <th class="text-left">Platform</th>
          <th class="text-left">Tags</th>
          <th class="text-left">Last Contacted</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in runners" :key="item.id">
          <td>
            <v-chip
              class="ma-2"
              :color="bgcolor(state(item))"
              :text-color="fgcolor(state(item))"
              label
              :to="{ name: 'RunnersView', params: { id: item.id } }"
            >
              <v-icon dark left>
                {{ icon(state(item)) }}
              </v-icon>
              {{ state(item) }}
            </v-chip>
          </td>
          <td>{{ item.name }}</td>
          <td>{{ item.description }}</td>
          <td>{{ item.platform }} / {{ item.architecture }}</td>
          <td>
            <v-chip
              class="ma-2"
              color="primary"
              small
              label
              v-if="item.run_untagged"
            >
              <v-icon left> mdi-label </v-icon>
              Run Untagged
            </v-chip>
            {{ item.tags }}
          </td>
          <td>
            {{ item.contacted_at }}
          </td>
        </tr>
      </tbody>
    </template>
  </v-simple-table>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "RegisterTokensList",

  created() {
    this.$store.dispatch("getRegisterTokens");
  },

  computed: mapState({
    tokens: (state) => state.tokes,
  }),
};
</script>
