<template>
  <div>
    <h1>Détails du Token</h1>
    <p>route: {{ $route.params.id }}</p>
    <p><strong>ID :</strong> {{ id }}</p>
    <p><strong>token :</strong> {{ name }}</p>
    <p><strong>supply :</strong> {{ supply }}</p>
    <p><strong>price :</strong> {{ price }}</p>
    <!-- Affichez d'autres informations sur le token ici si nécessaire -->
  </div>
</template>

<script>

export default {
  props: {
    id: {
      type: String,
      required: true
    },
    name: {
      type: String,
      required: true
    },
    price: {
      type: Number,
      required: true
    },
    supply: {
      type: Number,
      required: true
    },

  },
  mounted() {
    this.fetchTokens();
   
    
  },
  methods: {
    fetchTokens() {
      const tokenId = this.$route.params.id; 
      fetch(`http://localhost:8080/tokens/${tokenId}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          this.tokens = data; 
          console.log(data); 
        })
        .catch((error) => {
          console.error("Il y a eu un problème avec la requête fetch :", error);
        });
    }
}
}
</script>
