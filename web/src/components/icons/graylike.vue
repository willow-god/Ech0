<template>
  <svg xmlns="http://www.w3.org/2000/svg" width="1em" height="1em" viewBox="0 0 24 24">
    <!-- Icon from Material Line Icons by Vjacheslav Trushkin - https://github.com/cyberalien/line-md/blob/master/license.txt -->
    <path
      fill="none"
      stroke="#888888"
      stroke-dasharray="80"
      stroke-dashoffset="80"
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M7 11l5 -8l3 1l-1 6h7v3l-3 7h-11h-4v-9h4v9"
    >
      <animate fill="freeze" attributeName="stroke-dashoffset" dur="0.6s" values="80;0" />
    </path>
  </svg>
</template>

<script lang="ts">
export default {
  name: 'LineMdThumbsUp',
}
</script>
