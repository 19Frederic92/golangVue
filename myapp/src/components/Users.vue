<template>
  <div>
    <h1>Utilisateurs</h1>
    <ul>
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
      newUserName: '',
      newAge: null,
    };
  },
  mounted() {
    this.fetchUsers();
  },
  methods: {
    fetchUsers() {
      fetch("http://localhost:8080/users")
        .then((response) => response.json())
        .then((data) => {
          this.users = data;
          console.log(data)
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
        this.newAge='';
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
  },
};
</script>
