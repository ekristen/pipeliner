export default {
  methods: {
    color(state) {
      if (state == "success") {
        return "success";
      } else if (state == "skipped") {
        return "grey";
      } else if (state == "manual" || state == "canceled") {
        return "black";
      } else if (state == "active") {
        return "primary"
      } else if (state == "offline") {
        return "secondary"
      }

      return state;
    },
    icon(state) {
      if (state == "success") {
        return "mdi-checkbox-marked-circle";
      } else if (state == "skipped") {
        return "mdi-debug-step-over";
      } else if (state == "failed") {
        return "mdi-minus-circle-outline";
      } else if (state == "running") {
        return "mdi-play-circle-outline";
      } else if (state == "pending") {
        return "mdi-dots-horizontal-circle-outline";
      } else if (state == "created") {
        return "mdi-plus-circle-outline";
      } else if (state == "manual") {
        return "mdi-cog-outline";
      } else if (state == "canceled") {
        return "mdi-close-circle-outline";
      } else if (state == "active") {
        return "mdi-antenna";
      } else if (state == "offline") {
        return "mdi-minus-circle";
      }

      return "mdi-help-circle-outline";
    },
  }
}