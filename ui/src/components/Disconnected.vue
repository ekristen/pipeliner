<template>
  <div>
    <v-overlay :value="disconnected"></v-overlay>

    <v-dialog v-model="disconnected" hide-overlay persistent width="300">
      <v-card color="grey" dark>
        <v-card-text class="pt-3">
          Reconnecting
          <v-progress-linear
            indeterminate
            color="white"
            class="mb-0"
          ></v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "Disconnected",

  data: () => ({ disconnected: false }),

  computed: {
    ...mapState({
      socket: (state) => state.socket,
      isConnected: (state) => state.socket.isConnected,
    }),
  },
  watch: {
    isConnected(val) {
      this.disconnected = !val;
    },
  },
};
</script>