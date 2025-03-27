<template>
  <div>
    <h1>Warhammer Forbidden</h1>
    <ul>
      <div class="tableContainer">
        <div class="tableHeader">
          <p>id</p>
          <p>Name</p>
          <p>Image</p>
          
        </div>

        <!-- Lien vers la page de détail du token -->
        <ul>
          <li v-for="entites in entite" :key="entites.id">
            <div class="tokenRow">
              <p>{{ entites.id }}</p>
              <p>{{ entites.name }}</p>
              <img :src="entites.image">
              <!--<img src="/src/assets/entite/xenos.jpeg">-->
              
              <RouterLink 
  :to="{ 
    name: 'entiteDetail', 
    params: { 
      id: entites.id, 
      name: entites.name, 
      supply: entites.image, 
      
    } 
  }"
>
  {{ entites.name }}
</RouterLink>
             
            </div>
          </li>
        </ul>
</div>

      <li v-for="user in users" :key="user.id">
        {{ user.name }} 
        {{ user.age }}
        <button @click="editUser(user.id)">Modifier</button>
        <button @click="deleteUser(user.id)">Supprimer</button>
      </li>
    </ul>
        <!-- Input pour le nom de l'utilisateur -->
        <input v-model="newUserName" placeholder="Ajouter un nom d'utilisateur" />

<!-- Input pour l'âge de l'utilisateur -->
<input v-model="newAge" placeholder="Ajouter un âge" type="number" />

<!-- Bouton pour ajouter l'utilisateur -->
<button @click="addUser">Ajouter Utilisateur</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      users: [],
      entite: [],
      newUserName: '',
      newAge: null,
    };
  },
  mounted() {
    this.fetchEntite();
    this.fetchUsers();
    this.setupWebSocket();
  },
  methods: {
    setupWebSocket() {
      const socket = new WebSocket('ws://localhost:8080/ws');
      socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        if (message.type === 'update') {
          console.log('Mise à jour reçue:', message.data);
          this.fetchEntite(); // Recharger les données des entités
        }
      };
    },
    fetchUsers() {
      fetch("http://localhost:8080/users")
        .then((response) => response.json())
        .then((data) => {
          this.users = data;
          console.log(data);
        });
    },
    addUser() {
      fetch("http://localhost:8080/users", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name: this.newUserName, age: parseInt(this.newAge) }),
      })
      .then(() => {
        this.newUserName = '';
        this.newAge = '';
        this.fetchUsers(); // Recharger la liste des utilisateurs après ajout
      });
    },
    editUser(id) {
      const newName = prompt("Modifier le nom de l'utilisateur :");
      const newAge = prompt("Modifier l'age de l'utilisateur :");
      if (newName && newAge) {
        fetch(`http://localhost:8080/users/${id}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ name: newName, age: parseInt(newAge) }),
        })
        .then(() => {
          this.fetchUsers(); // Recharger la liste des utilisateurs après modification
        });
      }
    },
    deleteUser(id) {
      if (confirm("Êtes-vous sûr de vouloir supprimer cet utilisateur ?")) {
        fetch(`http://localhost:8080/users/${id}`, {
          method: 'DELETE',
        })
        .then(() => {
          this.fetchUsers(); // Recharger la liste des utilisateurs après suppression
        });
      }
    },
    fetchEntite() {
      console.log('Récupération des entités...');
      fetch("http://localhost:8080/entites")
        .then((response) => response.json())
        .then((data) => {
          this.entite = data;
          console.log('Données des entités mises à jour:', data);
        })
        .catch((error) => {
          console.error('Erreur lors de la récupération des entités:', error);
        });
    },
    fetchTokens() {
      fetch("http://localhost:8080/tokens")
        .then((response) => response.json())
        .then((data) => {
          this.tokens = data;
          console.log(data);
        });
    },
  },
};
</script>




<style>

.tableContainer {
  width: 100%;
  max-width: 600px; /* Largeur maximale pour limiter la taille */
  margin: auto;
}

.tableHeader, .tokenRow {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  font-weight: bold;
  background-color: #f1f1f1;
  border-bottom: 2px solid #ddd;
}

.tableHeader p {
  flex: 1;
  text-align: left;
}

.tokenRow {
  background-color: #ffffff;
  font-weight: normal;
}

.tokenRow p {
  flex: 1;
  text-align: left;
  margin: 0;
}

li {
  list-style-type: none;
}

ul {
  padding: 0;
  margin: 0;
}

li:not(:last-child) .tokenRow {
  border-bottom: 1px solid #ddd;
}

</style>
