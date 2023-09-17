
<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <h1>Search Flight Deals</h1>
        <div id="searchloc" class="container">
          <input v-model="searchQuery" @keyup.enter="fetchData" placeholder="Enter a location" />
          <button class="btn btn-primary" @click="fetchData">Search</button>
          <button class="btn btn-warning" @click="clearData">Clear</button>
          <br/>
          <br/>
          <div>
            <div id="dataBlock" class="vstack gap-3" v-if="showData">
              <p class="p-2" v-for="msg in apiData" :key="msg.id"><span v-html="msg.text"></span></p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      searchQuery: '',
      apiData: [],
      showData: false
    };
  },
  methods: {
    fetchData() {
      // Replace 'YOUR_API_URL' with the actual API endpoint, and pass the 'searchQuery' as a query parameter
      this.showData = true;
      const apiUrl = `http://localhost:8888/flights/${this.searchQuery}`;

      // Make the API request using Axios or Fetch
      axios
        .get(apiUrl)
        .then((response) => {
          this.apiData = response.data.data;
        })
        .catch((error) => {
          console.error(error);
        });
    },

    clearData() {
      this.showData = false;
      this.searchQuery = "";
      this.apiData = []
    }
  },
};
</script>
