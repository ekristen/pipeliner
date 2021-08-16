<template>
  <span>
    <v-chip
      v-for="tag in tags"
      :key="tag.tag"
      v-bind="$attrs"
      class="mr-2 mb-2 primary"
      label
    >{{ tag.tag }}</v-chip>
  </span>
</template>

<script>
export default {
  name: "Tags",

  props: ["tags"],
};
</script>