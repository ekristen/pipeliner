<template>
  <v-dialog v-model="dialog" max-width="600px" transition="dialog-bottom-transition">
    <template v-slot:activator="{ on, attrs }">
      <v-btn fab x-small color="accent" v-bind="attrs" v-on="on">
        <v-icon>mdi-pencil</v-icon>
      </v-btn>
    </template>

    <v-card>
      <v-toolbar color="primary" dark>Edit Runner: {{runner.name}}</v-toolbar>
      <v-card-text class="mt-3">
        <v-textarea
          label="Description"
          hint="Friendly description of the runner"
          v-model="description"
          auto-grow
          rows="1"
        ></v-textarea>
        <v-divider class="mt-2 mb-2"></v-divider>
        <p>Tags allow you to specific specific runners for specific jobs</p>
        <v-chip-group multiple color="tags" v-if="tags.length > 0">
          <v-chip v-for="tag in tags" :key="tag" close @click:close="removeTag(tag)">{{ tag }}</v-chip>
        </v-chip-group>
        <v-text-field
          label="Add Tags"
          hint="Type tags, seperated by a comma to add to the list"
          v-model="input"
        ></v-text-field>
        <v-divider class="mt-2 mb-2"></v-divider>
        <p>Run Untagged allows a runner to run a job if it has no tags</p>
        <v-switch v-model="run_untagged" label="Run Untagged" color="red darken-3"></v-switch>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="red darken-1" text @click="reset">Close</v-btn>
        <v-btn color="blue darken-1" text @click="save">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<script>
export default {
  name: "RunnerEdit",

  props: ["runner"],

  data: () => ({
    dialog: false,
    run_untagged: false,
    tags: [],
    input: "",
    description: "",
  }),

  created() {
    this.description = this.runner.description;
    this.run_untagged = this.runner.run_untagged;
    if (this.runner.tags) {
      this.tags = this.runner.tags.map((t) => t.tag);
    }
  },

  watch: {
    input(value) {
      if (value.indexOf(",") === value.length - 1) {
        this.addTag(value.split(",")[0]);
        this.input = "";
      }
    },
  },

  methods: {
    addTag(tag) {
      let newTag = tag.toLowerCase();
      if (newTag === "") {
        return;
      }
      if (!this.tags.includes(newTag)) {
        this.tags.push(newTag);
      }
    },
    removeTag(tag) {
      let newTag = tag.toLowerCase();
      this.tags = this.tags.filter((t) => t !== newTag);
    },
    reset() {
      this.dialog = false;
    },
    save() {
      this.$store.dispatch("saveRunner", {
        id: this.runner.id,
        description: this.description,
        tags: this.tags,
        run_untagged: this.run_untagged,
      });

      this.reset();
    },
  },
};
</script>