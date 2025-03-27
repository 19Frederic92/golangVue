<template>
  <div>
    <h1>Détails de l'entité</h1>
    <p>Route: {{ $route.params.id }}</p>
    <p><strong>ID :</strong> {{ entite.id }}</p>
    <p><strong>Nom :</strong> {{ entite.name }}</p>
    <p><strong>Image :</strong> {{ entite.image }}</p>
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
    image: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      entite: {
        id: this.id,
        name: this.name,
        image: this.image
      }
    };
  },
  mounted() {
    this.fetchEntite();
  },
  methods: {
    fetchEntite() {
      const entiteId = this.id;
      fetch(`http://localhost:8080/entites/${entiteId}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok");
          }
          return response.json();
        })
        .then((data) => {
          this.entite = data;
          console.log(data);
        })
        .catch((error) => {
          console.error("Il y a eu un problème avec la requête fetch :", error);
        });
    }
  }
};
</script>
