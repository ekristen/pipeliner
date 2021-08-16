<template>
  <v-form>
    <v-container fluid>
      <v-text-field
        label="Name"
        single-line
        full-width
        hide-details
        v-model="name"
      ></v-text-field>
    </v-container>
    <v-container fluid>
      <v-textarea
        label="Workflow Contents"
        counter
        auto-grow
        filled
        v-model="workflow"
      ></v-textarea>

      <v-btn class="mr-4" @click="submit"> Create </v-btn>
      <v-btn @click="clear"> clear </v-btn>
    </v-container>
  </v-form>
</template>

<script>
import Vue from "vue";

export default {
  name: "Run",

  data() {},

  methods: {
    submit() {},

    clear() {},
  },
};
</script>
