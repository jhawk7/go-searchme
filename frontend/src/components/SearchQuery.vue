
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
              <p class="p-2" v-if="noData">No Data</p>
              <p class="p-2" v-for="msg in paginatedData" :key="msg.id"><span v-html="msg.text"></span></p>
            </div>
          </div>
          <div class="pagination">
            <button class="btn btn-link" @click="previousPage" :disabled="currentPage === 1">Previous</button>
            <button class="btn btn-link" @click="nextPage" :disabled="currentPage === totalPages || apiData.length === 0">Next</button>
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
      offset: '',
      showData: false,
      currentPage: 1,
      itemsPerPage: 20,
    };
  },
  computed: {
    totalPages() {
      var totalItems = this.apiData.length;
      return Math.ceil(totalItems / this.itemsPerPage);
    },
    paginatedData() {
      var end = this.currentPage * this.itemsPerPage;
      var start = end - this.itemsPerPage;
      return this.apiData.slice(start,end);
    },
    noData() {
      return this.apiData.length === 0;
    }
  },
  methods: {
    fetchData() {
      // Replace 'YOUR_API_URL' with the actual API endpoint, and pass the 'searchQuery' as a query parameter
      this.showData = true;
      //const apiUrl = `/flights/${this.searchQuery}?offset=${this.offset}`;
      const apiUrl = `/v1/deals?filter=${this.searchQuery}&page=${this.currentPage}&pageSize=${this.itemsPerPage}`

      // Make the API request using Axios or Fetch
      axios
        .get(apiUrl)
        .then((response) => {
          this.apiData = response.data.messages;
          this.offset = response.data.offset;
        })
        .catch((error) => {
          console.error(error);
        });
    },

    previousPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },

    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },

    clearData() {
      this.showData = false;
      this.searchQuery = "";
      this.apiData = [];
      this.offset = "";
      this.currentPage = 1;
    }
  },
};
</script>

<style>
.pagination {
  justify-content: center;
}
</style>
